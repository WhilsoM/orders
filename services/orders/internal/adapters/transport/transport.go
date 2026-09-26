package transport

import (
	"backend/services/orders/internal/order"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type orderService interface {
	CreateOrder(ctx context.Context, items order.CreateOrderInput) (order.OrderCreated, error)
	GetOrder(ctx context.Context, orderID int) (order.Order, error)
}

type orderTransport struct {
	svc orderService
}

func NewOrderTransport(svc orderService) *orderTransport {
	return &orderTransport{
		svc,
	}
}

func (o *orderTransport) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if apiErr := readJSON(w, r, &req); apiErr != nil {
		slog.Info("create order invalid request", "status", apiErr.Status, "error", apiErr.Message)
		writeError(w, apiErr.Status, apiErr.Message)
		return
	}

	if err := req.Validate(); err != nil {
		slog.Info("create order failed validation", "error", err.Error())
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	items := make([]order.CreateOrderItemInput, len(req.Items))

	for i, item := range req.Items {
		items[i] = order.CreateOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	input := order.CreateOrderInput{
		Items: items,
	}

	order, err := o.svc.CreateOrder(r.Context(), input)
	if err != nil {
		slog.Error("failed to create order", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, CreateOrderResponse{
		OrderID: order.OrderID,
		Status:  order.Status,
	})
}

func (o *orderTransport) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		slog.Info("get order by id", "error", "id cannot be empty", "idStr", idStr)
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Info("get order by id", "error", "id cannot be not number", "idStr", idStr)
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	orderData, err := o.svc.GetOrder(r.Context(), id)
	if err != nil {
		slog.Info("get order by id", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	items := make([]GetOrderByIDItem, len(orderData.Items))

	for i, v := range orderData.Items {
		items[i] = GetOrderByIDItem{
			ID:        v.ID,
			ProductID: v.ProductID,
			Quantity:  v.Quantity,
			UnitPrice: v.UnitPrice,
		}
	}

	res := GetOrderByIDResponse{
		ID:          orderData.ID,
		Status:      orderData.Status,
		TotalAmount: orderData.TotalAmount,
		CreatedAt:   orderData.CreatedAt,
		UpdatedAt:   orderData.UpdatedAt,
		Items:       items,
	}

	writeJSON(w, http.StatusOK, res)
}
