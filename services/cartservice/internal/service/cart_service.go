package service

import (
	"bytes"
	"ecommerce-backend/services/cartservice/internal/model"
	"ecommerce-backend/services/cartservice/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

type CartService struct {
	Repo        *repository.CartRepository
	OrderSvcURL string // base URL of Order Service
}

func NewCartService(repo *repository.CartRepository, orderSvcURL string) *CartService {
	return &CartService{Repo: repo, OrderSvcURL: orderSvcURL}
}

func (s *CartService) AddItem(userID, productID string, qty int) (*model.Cart, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Get existing cart
	cart, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// If no cart, create one
	if cart == nil {
		cart = &model.Cart{
			UserID: userID,
			Items:  []model.CartItem{},
		}
		if err := s.Repo.Create(cart); err != nil {
			return nil, err
		}
	}

	// Check if product already exists in cart
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			cart.Items[i].Quantity += qty
			return cart, s.Repo.Save(cart)
		}
	}

	// Add new product
	cart.Items = append(cart.Items, model.CartItem{
		ProductID: productID,
		Quantity:  qty,
	})

	if err := s.Repo.Save(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// GetCart returns the current user's cart with items
func (s *CartService) GetCart(userID string) (*model.Cart, error) {
	return s.Repo.GetByUserID(userID)
}

// UpdateItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateItemQuantity(itemID string, qty int) (*model.CartItem, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	item, err := s.Repo.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}

	item.Quantity = qty
	if err := s.Repo.UpdateItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteItem removes a cart item from the cart
func (s *CartService) DeleteItem(itemID string) error {
	item, err := s.Repo.GetItemByID(itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}

	return s.Repo.DeleteItem(itemID)
}

func (s *CartService) Checkout(c *gin.Context, userID string) (map[string]interface{}, error) {
	cart, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")

	totalPrice := 0.0

	for i, item := range cart.Items {
		// Fetch product details from ProductService
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/products/%s", productServiceURL, item.ProductID))
		if err != nil {
			log.Printf("⚠️ Failed to fetch product %s: %v", item.ProductID, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("⚠️ Product %s not found in ProductService", item.ProductID)
			continue
		}

		var product struct {
			Product struct {
				ID    string  `json:"id"`
				Name  string  `json:"name"`
				Price float64 `json:"price"`
			} `json:"product"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			resp.Body.Close()
			log.Printf("⚠️ Error decoding product %s: %v", item.ProductID, err)
			continue
		}
		resp.Body.Close()

		price := product.Product.Price
		cart.Items[i].Price = price
		totalPrice += price * float64(item.Quantity)
	}

	// Prepare order payload
	orderPayload := map[string]interface{}{
		"user_id":     userID,
		"items":       cart.Items,
		"status":      "pending",
		"total_price": totalPrice,
	}

	body, _ := json.Marshal(orderPayload)

	// Call OrderService
	req, err := http.NewRequest("POST", s.OrderSvcURL+"/api/v1/orders", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Forward JWT if available
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call order service: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("🟢 OrderService Response:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create order in order service")
	}

	var orderResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		return nil, err
	}

	// ✅ Clear cart after successful order
	if err := s.Repo.ClearCart(cart.ID); err != nil {
		return nil, err
	}

	fmt.Printf("✅ Checkout completed. Total: ₹%.2f\n", totalPrice)
	return orderResp, nil
}
