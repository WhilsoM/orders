package main

import (
	"backend/services/orders/internal/adapters/transport"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	port := ":3000"

	r.Use(middleware.Logger)
	handler := transport.NewOrderTransport()

	r.Post("/orders", handler.CreateOrder)

	slog.Info("orders service is started", "port", port)

	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("server shutdowned", "error", err)
	}

}
