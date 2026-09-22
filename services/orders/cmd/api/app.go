package main

import (
	"backend/services/orders/internal/adapters/postgresql"
	"backend/services/orders/internal/adapters/transport"
	"backend/services/orders/internal/order"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newApp() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	repo := postgresql.NewOrdersRepository()
	svc := order.NewOrderService(repo)
	handler := transport.NewOrderTransport(svc)

	r.Post("/orders", handler.CreateOrder)
	r.Get("/orders/{id}", handler.GetOrderByID)

	return r
}
