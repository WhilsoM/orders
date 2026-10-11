package main

import (
	"backend/services/support/internal/adapters/transport"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	handler := &transport.Handler{}

	r.Get("/ws", handler.Connect)

	port := ":3001"

	slog.Info("support service is started", "port", port)

	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}
