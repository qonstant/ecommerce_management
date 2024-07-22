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

type OrdersHandler struct {
	db *postgres.Queries
}

func NewOrderHandler(conn *sql.DB) *OrdersHandler {
	return &OrdersHandler{
		db: postgres.New(conn),
	}
}

func (h *OrdersHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Get("/search/user", h.searchByUser)
	r.Get("/search/status", h.searchByStatus)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// @Summary List all orders
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} postgres.Order
// @Failure 500 {object} response.Object
// @Router /orders [get]
func (h *OrdersHandler) list(w http.ResponseWriter, r *http.Request) {
	orders, err := h.db.ListOrders(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, orders)
}

// @Summary Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body postgres.CreateOrderParams true "Order details"
// @Success 200 {object} postgres.Order
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders [post]
func (h *OrdersHandler) add(w http.ResponseWriter, r *http.Request) {
	var req postgres.CreateOrderParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	order, err := h.db.CreateOrder(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, order)
}

// @Summary Get an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} postgres.Order
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [get]
func (h *OrdersHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	order, err := h.db.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, order)
}

// @Summary Update an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body postgres.UpdateOrderParams true "Order details"
// @Success 200 {object} postgres.Order
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [put]
func (h *OrdersHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req postgres.UpdateOrderParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	order, err := h.db.UpdateOrder(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, order)
}

// @Summary Delete an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 204 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders/{id} [delete]
func (h *OrdersHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteOrder(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}

// @Summary Search orders by user ID
// @Tags orders
// @Accept json
// @Produce json
// @Param user query int true "User ID"
// @Success 200 {array} postgres.Order
// @Failure 500 {object} response.Object
// @Router /orders/search/user [get]
func (h *OrdersHandler) searchByUser(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.db.SearchOrdersByUser(r.Context(), userID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, orders)
}

// @Summary Search orders by status
// @Tags orders
// @Accept json
// @Produce json
// @Param status query string true "Order Status"
// @Success 200 {array} postgres.Order
// @Failure 500 {object} response.Object
// @Router /orders/search/status [get]
func (h *OrdersHandler) searchByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		response.BadRequest(w, r, errors.New("missing status parameter"), nil)
		return
	}

	orders, err := h.db.SearchOrdersByStatus(r.Context(), postgres.OrderStatus(status))
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, orders)
}
