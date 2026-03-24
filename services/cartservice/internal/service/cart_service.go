package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"ecommerce-backend/services/cartservice/internal/models"
	"ecommerce-backend/services/cartservice/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartService struct {
	Repo        *repository.CartRepository
	OrderSvcURL string
	HTTPClient  *http.Client
}

func NewCartService(repo *repository.CartRepository, orderSvcURL string) *CartService {
	return &CartService{
		Repo:        repo,
		OrderSvcURL: orderSvcURL,
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

type productResponse struct {
	Product struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Category string  `json:"category"`
		Price    float64 `json:"price"`
		Stock    int     `json:"stock"`
	} `json:"product"`
}

func (s *CartService) AddItem(userID, productID string, qty int) (*models.Cart, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	cart, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	product, err := s.fetchProduct(productID)
	if err != nil {
		return nil, err
	}

	productUUID, err := uuid.Parse(product.Product.ID)
	if err != nil {
		return nil, errors.New("invalid product ID in product service response")
	}

	existingQty := 0
	if cart != nil {
		for _, item := range cart.Items {
			if item.ProductID == productUUID {
				existingQty = item.Quantity
				break
			}
		}
	}

	if existingQty+qty > product.Product.Stock {
		return nil, fmt.Errorf(
			"not enough stock for product %s: available %d, requested %d",
			product.Product.ID,
			product.Product.Stock,
			existingQty+qty,
		)
	}

	if cart == nil {
		cart = &models.Cart{
			UserID: userUUID,
			Items:  []models.CartItem{},
		}
		if err := s.Repo.Create(cart); err != nil {
			return nil, err
		}
	}

	productDetails := &models.Product{
		ID:       productUUID,
		Name:     product.Product.Name,
		Price:    product.Product.Price,
		Stock:    product.Product.Stock,
		Category: product.Product.Category,
	}

	for i := range cart.Items {
		if cart.Items[i].ProductID == productUUID {
			cart.Items[i].Quantity += qty
			cart.Items[i].Price = product.Product.Price
			cart.Items[i].Product = productDetails
			return cart, s.Repo.Save(cart)
		}
	}

	cart.Items = append(cart.Items, models.CartItem{
		CartID:    cart.ID,
		ProductID: productUUID,
		Quantity:  qty,
		Price:     product.Product.Price,
		Product:   productDetails,
	})

	if err := s.Repo.Save(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) GetCart(userID string) (*models.Cart, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.Repo.GetByUserID(userID)
}

func (s *CartService) UpdateItemQuantity(itemID string, qty int) (*models.CartItem, error) {
	if qty <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	if _, err := uuid.Parse(itemID); err != nil {
		return nil, errors.New("invalid item id")
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

func (s *CartService) DeleteItem(itemID string) error {
	if _, err := uuid.Parse(itemID); err != nil {
		return errors.New("invalid item id")
	}

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
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	cart, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	totalPrice := 0.0
	for i, item := range cart.Items {
		product, err := s.fetchProduct(item.ProductID.String())
		if err != nil {
			return nil, err
		}

		cart.Items[i].Price = product.Product.Price
		cart.Items[i].Product = &models.Product{
			ID:       item.ProductID,
			Name:     product.Product.Name,
			Category: product.Product.Category,
			Price:    product.Product.Price,
			Stock:    product.Product.Stock,
		}
		totalPrice += product.Product.Price * float64(item.Quantity)
	}

	orderPayload := map[string]interface{}{
		"user_id":     userID,
		"items":       cart.Items,
		"status":      "pending",
		"total_price": totalPrice,
	}

	body, err := json.Marshal(orderPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.OrderSvcURL+"/api/v1/orders", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call order service: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create order in order service: %s", string(bodyBytes))
	}

	var orderResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		return nil, err
	}

	if err := s.Repo.ClearCart(cart.ID); err != nil {
		return nil, err
	}

	return orderResp, nil
}

func (s *CartService) CalculateTotal(items []models.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}

func (s *CartService) fetchProduct(productID string) (*productResponse, error) {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		return nil, errors.New("PRODUCT_SERVICE_URL not configured")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/products/%s", productServiceURL, productID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product %s not found in product service", productID)
	}

	var product productResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("invalid product response: %w", err)
	}
	if product.Product.ID == "" {
		return nil, errors.New("unexpected product payload")
	}

	return &product, nil
}
