package service

import (
	model "ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"errors"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

type OrderItemInput struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // in a real app, fetch from ProductService
}

func (s *OrderService) CreateOrder(userID string, items []OrderItemInput) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	var total float64
	orderItems := make([]model.OrderItem, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		total += float64(item.Quantity) * item.Price
		orderItems[i] = model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	order := &model.Order{
		UserID:     userID,
		Status:     "pending",
		TotalPrice: total,
		Items:      orderItems,
	}

	if err := s.Repo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

// List all orders for a user
func (s *OrderService) GetOrdersByUser(userID string) ([]model.Order, error) {
	return s.Repo.GetByUser(userID)
}

// Get single order by ID
func (s *OrderService) GetOrderByID(orderID string) (*model.Order, error) {
	return s.Repo.GetByID(orderID)
}
