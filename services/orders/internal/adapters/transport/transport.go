package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type orderTransport struct {
}

func NewOrderTransport() *orderTransport {
	return &orderTransport{}
}

func (o *orderTransport) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&req)

	slog.Info("create order", "product_id", req.ProductID, "quantity", req.Quantity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "received"}`))
}
