package payment

type CreatePaymentParams struct {
	UserID     int64   `json:"user_id"`
	OrderID    int64   `json:"order_id"`
	Amount     string  `json:"amount"`
	HPAN       string `json:"hpan"`
	ExpDate    string `json:"expDate"`
	CVC        string `json:"cvc"`
}