package postgresql

import (
	"backend/services/orders/internal/order"
	"errors"
	"log/slog"
)

type orderRepository struct {
	memory map[int][]order.Data
}

func NewOrdersRepository() *orderRepository {
	return &orderRepository{
		memory: map[int][]order.Data{
			1: []order.Data{
				{
					ProductID: 25,
					Quantity:  2,
				},
				{
					ProductID: 26,
					Quantity:  3,
				},
			},
		},
	}
}

func (o *orderRepository) CreateOrder(input order.CreateOrderInput) (string, error) {
	for _, val := range o.memory {
		for _, data := range val {
			for _, v := range input.Items {
				if data.ProductID == v.ProductID {
					slog.Info("found an item from repository",
						"product_id", data.ProductID,
						"quantity", data.Quantity,
					)
				}
			}
		}
	}

	return "received", nil
}

func (o *orderRepository) GetOrder(orderID int) ([]order.Data, error) {
	if val, ok := o.memory[orderID]; ok {
		slog.Info("get order by id", "order_id", orderID)
		return val, nil
	}

	return []order.Data{}, errors.New("not found an order")
}
