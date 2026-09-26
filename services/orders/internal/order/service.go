package order

import (
	"fmt"
	"log/slog"
)

type orderRepository interface {
	CreateOrder(input CreateOrderInput) (OrderCreated, error)
	GetOrder(orderID int) (Order, error)
}

type orderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *orderService {
	return &orderService{
		repo,
	}
}

func (o *orderService) CreateOrder(input CreateOrderInput) (OrderCreated, error) {
	slog.Info(
		"order created",
		"items_count", len(input.Items),
	)

	order, err := o.repo.CreateOrder(input)
	if err != nil {
		return OrderCreated{}, err
	}

	return order, nil
}

func (o *orderService) GetOrder(orderID int) (Order, error) {
	result, err := o.repo.GetOrder(orderID)
	if err != nil {
		return Order{}, fmt.Errorf("get order by id has an error: %w", err)
	}

	totalAmount := 0

	for _, v := range result.Items {
		totalAmount += int(v.UnitPrice) * v.Quantity
	}

	slog.Info("order", "total_amount", totalAmount)

	return Order{
		ID:          result.ID,
		Status:      result.Status,
		TotalAmount: totalAmount,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		Items:       result.Items,
	}, nil
}
