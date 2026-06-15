package events

type CheckoutRequested struct {
	UserID     string         `json:"user_id"`
	Items      []CheckoutItem `json:"items"`
	TotalPrice float64        `json:"total_price"`
}

type CheckoutItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
