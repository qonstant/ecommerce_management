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

type ProductsHandler struct {
	db *postgres.Queries
}

func NewProductHandler(conn *sql.DB) *ProductsHandler {
	return &ProductsHandler{
		db: postgres.New(conn),
	}
}

func (h *ProductsHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Get("/search/name", h.searchByName)
	r.Get("/search/category", h.searchByCategory)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// @Summary List all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} postgres.Product
// @Failure 500 {object} response.Object
// @Router /products [get]
func (h *ProductsHandler) list(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.ListProducts(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, products)
}

// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body postgres.CreateProductParams true "Product details"
// @Success 200 {object} postgres.Product
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /products [post]
func (h *ProductsHandler) add(w http.ResponseWriter, r *http.Request) {
	var req postgres.CreateProductParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	product, err := h.db.CreateProduct(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, product)
}

// @Summary Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} postgres.Product
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /products/{id} [get]
func (h *ProductsHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	product, err := h.db.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, product)
}

// @Summary Update a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body postgres.UpdateProductParams true "Product details"
// @Success 200 {object} postgres.Product
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /products/{id} [put]
func (h *ProductsHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req postgres.UpdateProductParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	product, err := h.db.UpdateProduct(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, product)
}

// @Summary Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /products/{id} [delete]
func (h *ProductsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}

// @Summary Search products by name
// @Tags products
// @Accept json
// @Produce json
// @Param name query string true "Product name"
// @Success 200 {array} postgres.Product
// @Failure 500 {object} response.Object
// @Router /products/search/name [get]
func (h *ProductsHandler) searchByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		response.BadRequest(w, r, errors.New("missing name parameter"), nil)
		return
	}

	// Convert string to sql.NullString
	searchName := sql.NullString{String: name, Valid: true}

	products, err := h.db.SearchProductsByName(r.Context(), searchName)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, products)
}


// @Summary Search products by category
// @Tags products
// @Accept json
// @Produce json
// @Param category query string true "Product category"
// @Success 200 {array} postgres.Product
// @Failure 500 {object} response.Object
// @Router /products/search/category [get]
func (h *ProductsHandler) searchByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		response.BadRequest(w, r, errors.New("missing category parameter"), nil)
		return
	}

	products, err := h.db.SearchProductsByCategory(r.Context(), category)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, products)
}
