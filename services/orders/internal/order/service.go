package order

import (
	"fmt"
	"log/slog"
)

type orderRepository interface {
	CreateOrder(input CreateOrderInput) (string, error)
	GetOrder(orderID int) ([]Data, error)
}

type orderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *orderService {
	return &orderService{
		repo,
	}
}

func (o *orderService) CreateOrder(input CreateOrderInput) (string, error) {
	slog.Info(
		"order created",
		"items_count", len(input.Items),
	)

	status, err := o.repo.CreateOrder(input)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (o *orderService) GetOrder(orderID int) ([]Data, error) {
	result, err := o.repo.GetOrder(orderID)
	if err != nil {
		return []Data{}, fmt.Errorf("get order by id has an error: %w", err)
	}

	return result, nil
}
