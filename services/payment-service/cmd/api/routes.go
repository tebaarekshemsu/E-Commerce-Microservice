package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	// Match Java context-path /payment-service
	mux.Route("/payment-service", func(r chi.Router) {
		r.Route("/api/payments", func(r chi.Router) {
			r.Get("/", app.GetAllPayments)
			r.Get("/{paymentId}", app.GetPayment)
			r.Post("/", app.CreatePayment)
			r.Put("/", app.UpdatePayment)
			r.Delete("/{paymentId}", app.DeletePayment)
		})
	})

	return mux
}
