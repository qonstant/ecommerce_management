package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ecommerce_management/internal/repository/postgres"
	"ecommerce_management/pkg/server/response"
)

type PaymentsHandler struct {
	db *postgres.Queries
}

func NewPaymentHandler(conn *sql.DB) *PaymentsHandler {
	return &PaymentsHandler{
		db: postgres.New(conn),
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

// @Summary Create a new payment
// @Tags payments
// @Accept json
// @Produce json
// @Param request body postgres.CreatePaymentParams true "Payment details"
// @Success 200 {object} postgres.Payment
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /payments [post]
func (h *PaymentsHandler) add(w http.ResponseWriter, r *http.Request) {
	var req postgres.CreatePaymentParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	payment, err := h.db.CreatePayment(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	// Here you would make a call to ePayment.kz API for actual payment processing
	// Ensure you handle that part appropriately

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
		response.BadRequest(w, r, errors.New("missing user parameter"), nil)
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
		response.BadRequest(w, r, errors.New("missing order parameter"), nil)
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
// @Param status query string true "Payment status"
// @Success 200 {array} postgres.Payment
// @Failure 500 {object} response.Object
// @Router /payments/search/status [get]
func (h *PaymentsHandler) searchByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		response.BadRequest(w, r, errors.New("missing status parameter"), nil)
		return
	}

	// Convert string to PaymentStatus
	paymentStatus := postgres.PaymentStatus(status)

	payments, err := h.db.SearchPaymentsByStatus(r.Context(), paymentStatus)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, payments)
}
