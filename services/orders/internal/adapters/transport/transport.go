package transport

import (
	"backend/services/orders/internal/order"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type orderServiceI interface {
	CreateOrder(items order.CreateOrderInput) (string, error)
	GetOrder(orderID int) ([]order.Data, error)
}

type orderTransport struct {
	svc orderServiceI
}

func NewOrderTransport(svc orderServiceI) *orderTransport {
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

	status, err := o.svc.CreateOrder(input)
	if err != nil {
		slog.Error("failed to create order", "error", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, CreateOrderResponse{
		Status: status,
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

	order, err := o.svc.GetOrder(id)
	if err != nil {
		slog.Info("get order by id", "error", err)
		writeError(w, http.StatusBadGateway, "something went wrong...")
		return
	}

	res := make([]GetOrderByIDResponse, len(order))

	for i, val := range order {
		res[i] = GetOrderByIDResponse{
			ProductID: val.ProductID,
			Quantity:  val.Quantity,
		}
	}

	writeJSON(w, http.StatusOK, res)
}
