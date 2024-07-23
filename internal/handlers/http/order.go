package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ecommerce_management/internal/repository/postgres"
	"ecommerce_management/internal/domain/order"
	"ecommerce_management/pkg/server/response"
)

type OrdersHandler struct {
	store *postgres.Store
}

func NewOrderHandler(db *sql.DB) *OrdersHandler {
	return &OrdersHandler{
		store: postgres.NewStore(db),
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
	orders, err := h.store.ListOrders(r.Context())
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
// @Param request body order.CreateOrderRequest true "Order details"
// @Success 200 {object} postgres.Order
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /orders [post]
func (h *OrdersHandler) add(w http.ResponseWriter, r *http.Request) {
	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	tx, err := h.store.BeginTx(r.Context(), nil)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	// Create the order first with a placeholder total amount ("0.00")
	order, err := tx.CreateOrder(r.Context(), postgres.CreateOrderParams{
		UserID:      req.UserID,
		TotalAmount: "0.00", // Dummy value, will be updated later
	})
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	var totalAmount float64
	for _, item := range req.Items {
		// Fetch the product
		product, err := tx.GetProduct(r.Context(), item.ProductID)
		if err != nil {
			return
		}

		// Check stock quantity
		if product.StockQuantity < item.Quantity {
			response.BadRequest(w, r, fmt.Errorf("insufficient stock for product ID %d", item.ProductID), nil)
			return
		}

		// Convert product.Price to float64
		productPrice, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			response.InternalServerError(w, r, fmt.Errorf("invalid product price format: %v", err))
			return
		}

		// Calculate item price
		itemPrice := productPrice * float64(item.Quantity)
		totalAmount += itemPrice

		// Format itemPrice to string with 2 decimal places
		itemPriceStr := fmt.Sprintf("%.2f", itemPrice)

		_, err = tx.CreateOrderItem(r.Context(), postgres.CreateOrderItemParams{
			OrderID:   order.ID, // Use the ID of the newly created order
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPriceStr, // Use the formatted string
		})
		if err != nil {
			response.InternalServerError(w, r, err)
			return
		}

		_, err = tx.UpdateProductStock(r.Context(), postgres.UpdateProductStockParams{
			StockQuantity: item.Quantity,
			ID:            item.ProductID,
		})
		if err != nil {
			response.InternalServerError(w, r, err)
			return
		}
	}

	// Format totalAmount to have at most 2 decimal places
	totalAmountStr := fmt.Sprintf("%.2f", totalAmount)

	// Update the order with the final totalAmount
	_, err = tx.UpdateOrder(r.Context(), postgres.UpdateOrderParams{
		ID:          order.ID, // Use the ID of the newly created order
		UserID:      order.UserID,
		TotalAmount: totalAmountStr, // Use the formatted string
		Status:      order.Status,
	})
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	// Re-fetch the updated order to return the correct total amount
	updatedOrder, err := h.store.GetOrder(r.Context(), order.ID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, updatedOrder)
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

	order, err := h.store.GetOrder(r.Context(), id)
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

	order, err := h.store.UpdateOrder(r.Context(), req)
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

	err = h.store.DeleteOrder(r.Context(), id)
	if err != nil {
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
// @Param user_id query int true "User ID"
// @Success 200 {array} postgres.Order
// @Failure 500 {object} response.Object
// @Router /orders/search/user [get]
func (h *OrdersHandler) searchByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	orders, err := h.store.SearchOrdersByUser(r.Context(), userID)
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
// @Param status query string true "Order status"
// @Success 200 {array} postgres.Order
// @Failure 500 {object} response.Object
// @Router /orders/search/status [get]
func (h *OrdersHandler) searchByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	orders, err := h.store.SearchOrdersByStatus(r.Context(), postgres.OrderStatus(status))
	
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, orders)
}

