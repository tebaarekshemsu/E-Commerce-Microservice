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

	mux.Route("/api/products", func(r chi.Router) {
		r.Get("/", app.GetAllProducts)
		r.Get("/{productId}", app.GetProduct)
		r.Post("/", app.CreateProduct)
		r.Put("/", app.UpdateProduct)
		r.Put("/{productId}", app.UpdateProductWithID)
		r.Delete("/{productId}", app.DeleteProduct)
	})

	mux.Route("/api/categories", func(r chi.Router) {
		r.Get("/", app.GetAllCategories)
		r.Get("/{categoryId}", app.GetCategory)
		r.Post("/", app.CreateCategory)
		r.Put("/", app.UpdateCategory)
		r.Put("/{categoryId}", app.UpdateCategoryWithID)
		r.Delete("/{categoryId}", app.DeleteCategory)
	})

	return mux
}
