package order

// CreateOrderRequest represents the request payload for creating a new order with items.
type CreateOrderRequest struct {
    UserID int64       `json:"user_id"`
    Items  []OrderItem `json:"items"`
}

// OrderItem represents an item in the order.
type OrderItem struct {
	ProductID int64 `json:"product_id"` // The ID of the product being ordered
	Quantity  int32 `json:"quantity"`   // The quantity of the product
}


