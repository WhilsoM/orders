package main

import (
	"backend/pkg/psql"
	"backend/services/orders/internal/adapters/postgresql"
	"backend/services/orders/internal/adapters/transport"
	"backend/services/orders/internal/order"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newApp() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	client := psql.PostgresClient{
		Host:     "localhost",
		User:     "postgres",
		Password: "mysecretpassword",
		Dbname:   "orders",
		Port:     "5432",
		Sslmode:  "disable",
	}
	db, err := psql.NewPostgres(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	repo := postgresql.NewOrdersRepository(db)
	svc := order.NewOrderService(repo)
	handler := transport.NewOrderTransport(svc)

	r.Post("/orders", handler.CreateOrder)
	r.Get("/orders/{id}", handler.GetOrderByID)

	return r
}
