package main

import (
	"backend/services/orders/internal/config"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.MustConfig()
	app := newApp(cfg)

	slog.Info("orders service is started", "port", cfg.HttpPort)

	if err := http.ListenAndServe(cfg.HttpPort, app); err != nil {
		panic(err)
	}
}
