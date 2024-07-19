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

type UsersHandler struct {
	db *postgres.Queries
}

func NewUserHandler(conn *sql.DB) *UsersHandler {
	return &UsersHandler{
		db: postgres.New(conn),
	}
}

func (h *UsersHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Get("/search/email", h.searchByEmail)
	r.Get("/search/name", h.searchByName)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// @Summary List of users from the repository
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} postgres.User
// @Failure 500 {object} response.Object
// @Router /users [get]
func (h *UsersHandler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, users)
}

// @Summary Add a new user to the repository
// @Tags users
// @Accept json
// @Produce json
// @Param request body postgres.CreateUserParams true "User details"
// @Success 200 {object} postgres.User
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /users [post]
func (h *UsersHandler) add(w http.ResponseWriter, r *http.Request) {
	var req postgres.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	user, err := h.db.CreateUser(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, user)
}

// @Summary Get a user from the repository
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} postgres.User
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /users/{id} [get]
func (h *UsersHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	user, err := h.db.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, user)
}

// @Summary Update a user in the repository
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body postgres.UpdateUserParams true "User details"
// @Success 200 {object} postgres.User
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /users/{id} [put]
func (h *UsersHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req postgres.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	user, err := h.db.UpdateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, user)
}

// @Summary Delete a user from the repository
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /users/{id} [delete]
func (h *UsersHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}

// @Summary Search users by email
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {array} postgres.User
// @Failure 500 {object} response.Object
// @Router /users/search/email [get]
func (h *UsersHandler) searchByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response.BadRequest(w, r, errors.New("missing email parameter"), nil)
		return
	}

	users, err := h.db.SearchUsersByEmail(r.Context(), email)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, users)
}

// @Summary Search users by name
// @Tags users
// @Accept json
// @Produce json
// @Param name query string true "User Name"
// @Success 200 {array} postgres.User
// @Failure 500 {object} response.Object
// @Router /users/search/name [get]
func (h *UsersHandler) searchByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		response.BadRequest(w, r, errors.New("missing name parameter"), nil)
		return
	}

	users, err := h.db.SearchUsersByName(r.Context(), sql.NullString{String: name, Valid: true})
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, users)
}
