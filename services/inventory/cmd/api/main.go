package main

import (
	"log/slog"
	"net/http"
)

func main() {
	app := newApp()

	port := ":3000"

	slog.Info("orders service is started", "port", port)

	if err := http.ListenAndServe(port, app); err != nil {
		panic(err)
	}
}
