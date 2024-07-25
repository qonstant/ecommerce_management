package epay

import "net/http"

type CreateInvoiceRequest struct {
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Name        string `json:"name"`
	Cryptogram  string `json:"cryptogram"`
	InvoiceID   string `json:"invoiceId"`
	Description string `json:"description"`
	Email       string `json:"email"`
	CardSave    bool   `json:"cardSave"`
	PostLink    string `json:"postLink"`
}

type CreateInvoiceResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type TokenResponse struct {
	Scope        string `json:"scope"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

type Cryptogram struct {
	HPAN       string `json:"hpan"`
	ExpDate    string `json:"expDate"`
	CVC        string `json:"cvc"`
	TerminalID string `json:"terminalId"`
}

type Credentials struct {
	URL            string
	Login          string
	Password       string
	OAuthURL       string
	PaymentPageURL string
	ShopID         string
	TerminalID     string
	GlobalToken    TokenResponse
}

type Client struct {
	httpClient  *http.Client
	Credentials Credentials
}
