package epay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func New(credentials Credentials) (*Client, error) {
	// Ensure that all required fields are provided
	if credentials.OAuthURL == "" {
		return nil, errors.New("OAuthURL cannot be empty")
	}
	if credentials.ShopID == "" {
		return nil, errors.New("ShopID cannot be empty")
	}
	if credentials.TerminalID == "" {
		return nil, errors.New("TerminalID cannot be empty")
	}

	// Create and return the new Client instance
	client := &Client{
		httpClient:  &http.Client{Timeout: 40 * time.Second}, // Initialize httpClient here
		Credentials: credentials,
	}

	if err := client.initGlobalTokenRefresher(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) request(ctx context.Context, repeat bool, method, url string, body io.Reader, headers map[string]string, out interface{}) (err error) {
	// setup http request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// setup request header
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// send http request
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer res.Body.Close()

	// check unauthorized status
	if res.StatusCode == http.StatusUnauthorized && repeat {
		if err = c.initGlobalTokenRefresher(); err != nil {
			fmt.Printf("Error refreshing token: %v\n", err)
			return
		}
		return c.request(ctx, false, method, url, body, headers, out)
	}

	// read response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// check error status
	if res.StatusCode != http.StatusOK {
		return errors.New(string(data))
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
	}

	return
}

// CreateInvoice creates a new invoice and returns the response.
func (c *Client) CreateInvoice(ctx context.Context, token string, req CreateInvoiceRequest) (CreateInvoiceResponse, error) {
	var response CreateInvoiceResponse

	// Define the endpoint URL for creating an invoice (adjust as necessary)
	url := c.Credentials.PaymentPageURL

	// Prepare headers
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	// Convert the request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return response, err
	}

	// Make the request
	err = c.request(ctx, true, http.MethodPost, url, io.NopCloser(bytes.NewReader(body)), headers, &response)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return response, err
	}

	response.Success = true

	return response, nil
}
