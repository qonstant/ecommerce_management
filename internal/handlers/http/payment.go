package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"ecommerce_management/internal/domain/payment"
	"ecommerce_management/internal/provider/epay"
	"ecommerce_management/internal/repository/postgres"
	"ecommerce_management/pkg/server/response"
	"fmt"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/rand"
)

type PaymentsHandler struct {
	db         *postgres.Queries
	epayClient *epay.Client
	shopID     string
	terminalID string
}

func NewPaymentsHandler(conn *sql.DB, epayClient *epay.Client) *PaymentsHandler {
	return &PaymentsHandler{
		db:         postgres.New(conn),
		epayClient: epayClient,
		shopID:     epayClient.Credentials.ShopID,
		terminalID: epayClient.Credentials.TerminalID,
	}
}

func (h *PaymentsHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	r.Get("/search/user", h.searchByUser)
	r.Get("/search/order", h.searchByOrder)
	r.Get("/search/status", h.searchByStatus)

	return r
}

// @Summary List all payments
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {array} postgres.Payment
// @Failure 500 {object} response.Object
// @Router /payments [get]
func (h *PaymentsHandler) list(w http.ResponseWriter, r *http.Request) {
	payments, err := h.db.ListPayments(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, payments)
}

func generateInvoiceID() string {
	rand.Seed(uint64(time.Now().UnixNano())) // Convert int64 to uint64
	return fmt.Sprintf("%012d", rand.Int63n(1e12))
}

// @Summary Create a new payment
// @Tags payments
// @Accept json
// @Produce json
// @Param request body payment.CreatePaymentParams true "Payment details"
// @Success 200 {object} postgres.Payment
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments [post]
// @Summary Create a new payment
// @Tags payments
// @Accept json
// @Produce json
// @Param request body payment.CreatePaymentParams true "Payment details"
// @Success 200 {object} postgres.Payment
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments [post]
func (h *PaymentsHandler) add(w http.ResponseWriter, r *http.Request) {
	var req payment.CreatePaymentParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	amount, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		response.BadRequest(w, r, fmt.Errorf("invalid amount format"), req)
		return
	}

	user, err := h.db.GetUser(r.Context(), req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(w, r, fmt.Errorf("user not found"))
			return
		}
		response.InternalServerError(w, r, err)
		return
	}

	payment, err := h.db.CreatePayment(r.Context(), postgres.CreatePaymentParams{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  postgres.PaymentStatusUnsuccessful,
	})
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	invoiceID := generateInvoiceID()

	src := epay.PaymentRequest{
		Amount:    req.Amount,
		Currency:  "KZT",
		InvoiceID: invoiceID,
	}

	token, err := h.epayClient.GetPaymentToken(r.Context(), &src)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	cryptogram := epay.Cryptogram{
		HPAN:       req.HPAN,
		ExpDate:    req.ExpDate,
		CVC:        req.CVC,
		TerminalID: h.terminalID,
	}

	jsonData, err := json.Marshal(cryptogram)
	if err != nil {
		log.Printf("Error marshaling cryptogram to JSON: %v", err)
		response.InternalServerError(w, r, err)
		return
	}

	cryptogramString, err := epay.EncryptWithPublicKey(jsonData, epay.PublicKeyPEM)
	if err != nil {
		log.Printf("Error encrypting data: %v", err)
		response.InternalServerError(w, r, err)
		return
	}

	invoiceReq := epay.CreateInvoiceRequest{
		Amount:      amount,
		Currency:    "KZT",
		Name:        user.FullName,
		Cryptogram:  cryptogramString,
		Email:       user.Email,
		InvoiceID:   invoiceID,
		Description: "Payment for Order " + fmt.Sprint(req.OrderID),
		CardSave:    false,
		PostLink:    "https://testmerchant/order/" + fmt.Sprint(req.OrderID),
	}

	invoiceResp, err := h.epayClient.CreateInvoice(r.Context(), token.AccessToken, invoiceReq)
	if err != nil {
		log.Printf("Error creating invoice: %v", err)
		response.InternalServerError(w, r, err)
		return
	}

	// Update the payment status based on the response from the invoice creation
	var status postgres.PaymentStatus
	if invoiceResp.Success {
		status = postgres.PaymentStatusSuccessful
	} else {
		status = postgres.PaymentStatusUnsuccessful
	}

	// Update the payment record in the database with the new status
	updateReq := postgres.UpdatePaymentParams{
		ID:      payment.ID,
		UserID:  payment.UserID,
		OrderID: payment.OrderID,
		Amount:  payment.Amount,
		Status:  status,
	}

	payment, err = h.db.UpdatePayment(r.Context(), updateReq)
	if err != nil {
		log.Printf("Error updating payment: %v", err)
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, payment)
}

// @Summary Get a payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} postgres.Payment
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [get]
func (h *PaymentsHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	payment, err := h.db.GetPayment(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, payment)
}

// @Summary Update a payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param request body postgres.UpdatePaymentParams true "Payment details"
// @Success 200 {object} postgres.Payment
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [put]
func (h *PaymentsHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req postgres.UpdatePaymentParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	payment, err := h.db.UpdatePayment(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, payment)
}

// @Summary Delete a payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 204 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments/{id} [delete]
func (h *PaymentsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeletePayment(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}

// @Summary Search payments by user ID
// @Tags payments
// @Accept json
// @Produce json
// @Param user query int true "User ID"
// @Success 200 {array} postgres.Payment
// @Failure 500 {object} response.Object
// @Router /payments/search/user [get]
func (h *PaymentsHandler) searchByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user")
	if userIDStr == "" {
		response.BadRequest(w, r, errors.New("user query parameter is required"), nil)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	payments, err := h.db.SearchPaymentsByUser(r.Context(), userID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, payments)
}

// @Summary Search payments by order ID
// @Tags payments
// @Accept json
// @Produce json
// @Param order query int true "Order ID"
// @Success 200 {array} postgres.Payment
// @Failure 500 {object} response.Object
// @Router /payments/search/order [get]
func (h *PaymentsHandler) searchByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order")
	if orderIDStr == "" {
		response.BadRequest(w, r, errors.New("order query parameter is required"), nil)
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	payments, err := h.db.SearchPaymentsByOrder(r.Context(), orderID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, payments)
}

// @Summary Search payments by status
// @Tags payments
// @Accept json
// @Produce json
// @Param status query string true "Payment Status"
// @Success 200 {array} postgres.Payment
// @Failure 500 {object} response.Object
// @Router /payments/search/status [get]
func (h *PaymentsHandler) searchByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		response.BadRequest(w, r, errors.New("status query parameter is required"), nil)
		return
	}

	payments, err := h.db.SearchPaymentsByStatus(r.Context(), postgres.PaymentStatus(status))
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, payments)
}
