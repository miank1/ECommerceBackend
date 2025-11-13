package service

import (
	"bytes"
	"ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

type OrderItemReq struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func (s *OrderService) CreateOrder(userID string, items []OrderItemReq) (*models.Order, error) {
	var total float64
	for _, i := range items {
		total += float64(i.Quantity) * i.Price
	}

	order := models.Order{
		UserID:     userID,
		Status:     "pending",
		TotalPrice: total,
	}

	for _, i := range items {
		order.Items = append(order.Items, models.OrderItem{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
			Price:     i.Price,
		})
	}

	if err := s.Repo.Create(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	return s.Repo.GetByID(id)
}

func (s *OrderService) UpdateStatus(orderID, status string) (*models.Order, error) {
	order, err := s.Repo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	order.Status = status

	if err := s.Repo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to update order status: %v", err)
	}

	return order, nil
}

func (s *OrderService) UpdateInventory(orderID string) error {
	order, err := s.Repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		return fmt.Errorf("PRODUCT_SERVICE_URL not configured")
	}

	for _, item := range order.Items {
		payload := map[string]interface{}{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
		}

		body, _ := json.Marshal(payload)

		resp, err := http.Post(fmt.Sprintf("%s/api/v1/products/reduce-stock", productServiceURL),
			"application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf(" Failed to update stock for %s: %v\n", item.ProductID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf(" Stock update failed for product %s (status %d)\n", item.ProductID, resp.StatusCode)
		} else {
			fmt.Printf(" Stock updated for product %s\n", item.ProductID)
		}
	}

	return nil
}
