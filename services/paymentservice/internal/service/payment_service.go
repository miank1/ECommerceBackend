package service

import (
	"bytes"
	"ecommerce-backend/services/paymentservice/internal/models"
	"ecommerce-backend/services/paymentservice/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	Repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{Repo: repo}
}

// ProcessPayment handles payment creation and updates the order if successful
func (s *PaymentService) ProcessPayment(orderID, userID string, amount float64, method string) (*models.Payment, error) {
	// Simulate random payment success/failure (for testing)
	status := "success"
	if rand.Intn(100) < 10 { // 10% failure rate
		status = "failed"
	}

	payment := &models.Payment{
		OrderID:       uuid.MustParse(orderID),
		UserID:        uuid.MustParse(userID),
		Amount:        amount,
		Status:        status,
		PaymentMethod: method,
		TransactionID: fmt.Sprintf("trxn_%d", time.Now().UnixNano()),
	}

	// Save payment to DB
	if err := s.Repo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to record payment: %w", err)
	}

	// Asynchronously notify Order Service
	if status == "success" {
		go func() {
			if err := s.NotifyOrderService(orderID, "paid"); err != nil {
				log.Println("❌ Failed to notify OrderService:", err)
			} else {
				log.Println("✅ OrderService notified: order marked as paid")
			}
		}()
	}

	return payment, nil
}

// NotifyOrderService tells the Order Service that an order has been paid
func (s *PaymentService) NotifyOrderService(orderID, newStatus string) error {
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		return fmt.Errorf("ORDER_SERVICE_URL not configured")
	}

	payload := map[string]string{
		"status": newStatus,
	}

	body, _ := json.Marshal(payload)

	// ✅ Correct endpoint
	url := fmt.Sprintf("%s/api/v1/orders/%s/update-status", orderServiceURL, orderID)

	log.Println("Order URL is -----------------", url)

	// Create request
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ✅ Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call OrderService: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OrderService returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
