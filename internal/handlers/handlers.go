package handlers

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	healthPg "github.com/hellofresh/health-go/v5/checks/postgres"
	httpSwagger "github.com/swaggo/http-swagger"

	"ecommerce_management/docs"
	"ecommerce_management/internal/config"
	"ecommerce_management/internal/handlers/http"
	"ecommerce_management/internal/provider/epay"
)

type Dependencies struct {
	DB         *sql.DB
	Configs    config.Config
	EpayClient *epay.Client
}

// Configuration is an alias for a function that modifies the Handler
type Configuration func(h *Handler) error

// Handler is an implementation of the Handler
type Handler struct {
	dependencies Dependencies
	HTTP         *chi.Mux
}

// New creates a new Handler
func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// Create the handler
	h = &Handler{
		dependencies: d,
	}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

// @title E-commerce Management Service
// @version 1.0
// @description This is a simple E-commerce Management Service

// WithHTTPHandler applies an HTTP handler to the Handler
func WithHTTPHandler() Configuration {
	return func(h *Handler) error {
		// Create the HTTP handler
		h.HTTP = chi.NewRouter()

		// Init swagger handler
		docs.SwaggerInfo.BasePath = h.dependencies.Configs.BaseURL
		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		// Init service handlers
		userHandler := http.NewUserHandler(h.dependencies.DB)
		productHandler := http.NewProductHandler(h.dependencies.DB)
		orderHandler := http.NewOrderHandler(h.dependencies.DB)
		paymentHandler := http.NewPaymentsHandler(h.dependencies.DB, h.dependencies.EpayClient)

		h.HTTP.Route("/", func(r chi.Router) {
			r.Mount("/users", userHandler.Routes())
			r.Mount("/products", productHandler.Routes())
			r.Mount("/orders", orderHandler.Routes())

			r.Mount("/payments", paymentHandler.Routes())
		})

		// Setting up health checks
		healthHandler, _ := health.New(health.WithComponent(health.Component{
			Name:    "ecommerce-management-service",
			Version: "v1.0",
		}), health.WithChecks(
			health.Config{
				Name:      "postgres",
				Timeout:   time.Second * 10,
				SkipOnErr: false,
				Check: healthPg.New(healthPg.Config{
					DSN: os.Getenv("DB_SOURCE"),
				}),
			},
		))

		// Registering health check endpoint
		h.HTTP.Get("/status", healthHandler.HandlerFunc)

		return nil
	}
}
