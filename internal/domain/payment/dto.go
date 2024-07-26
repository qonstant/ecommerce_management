package payment

type CreatePaymentParams struct {
	OrderID    int64   `json:"order_id"`
	HPAN       string `json:"hpan"`
	ExpDate    string `json:"expDate"`
	CVC        string `json:"cvc"`
}