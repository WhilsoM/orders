package main

import (
	"backend/pkg/psql"
	"backend/services/orders/internal/adapters/inventorygrpc"
	"backend/services/orders/internal/adapters/postgresql"
	"backend/services/orders/internal/adapters/transport"
	"backend/services/orders/internal/config"
	"backend/services/orders/internal/order"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newApp(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	client := psql.PostgresClient{
		Host:            cfg.DBHost,
		User:            cfg.DBUser,
		Password:        cfg.DBPassword,
		Dbname:          cfg.DBDbname,
		Port:            cfg.DBPort,
		Sslmode:         cfg.DBSslmode,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime * time.Minute,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime * time.Minute,
	}
	db, err := psql.NewPostgres(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	repo := postgresql.NewOrdersRepository(db)

	inventoryClient, err := inventorygrpc.NewClient(cfg.GrpcInventoryClient)
	if err != nil {
		log.Fatal(err)
	}

	svc := order.NewOrderService(repo, inventoryClient)
	handler := transport.NewOrderTransport(svc)

	r.Post("/orders", handler.CreateOrder)
	r.Get("/orders/{id}", handler.GetOrderByID)

	return r
}
