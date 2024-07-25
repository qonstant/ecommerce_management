package epay

import "net/http"

type CreateInvoiceRequest struct {
	ShopID     string `json:"shop_id"`
	TerminalID string `json:"terminal_id"`
	OrderID    string `json:"order_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
}

type CreateInvoiceResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type TokenResponse struct {
	Scope        string `json:"scope"`
	ExpiresIn    string `json:"expires_in"` // Changed to string
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

type Credentials struct {
	URL            string
	Login          string
	Password       string
	OAuthURL       string
	PaymentPageURL string
	ShopID          string
	TerminalID      string
	GlobalToken    TokenResponse
}

type Client struct {
	httpClient  *http.Client
	credentials Credentials
}
