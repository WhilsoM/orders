package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type InventoryClient interface {
	GetItems(ctx context.Context, productIDs []int64) ([]InventoryItem, error)
}

type orderRepository interface {
	CreateOrder(input []InventoryItem) (OrderCreated, error)
	GetOrder(orderID int) (Order, error)
}

type orderService struct {
	repo      orderRepository
	inventory InventoryClient
}

func NewOrderService(repo orderRepository, inventory InventoryClient) *orderService {
	return &orderService{
		repo:      repo,
		inventory: inventory,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, input CreateOrderInput) (OrderCreated, error) {
	slog.Info(
		"create order",
		"items_count", len(input.Items),
	)
	// remove duplicates
	aggregatedItems := make(map[int]Data, 0)

	slog.Info("inventory", "items", o.inventory)

	for _, val := range input.Items {
		item := aggregatedItems[val.ProductID]
		item.Quantity += val.Quantity

		aggregatedItems[val.ProductID] = item
	}

	// get only ids
	productIDs := make([]int64, 0)

	for productID := range aggregatedItems {
		productIDs = append(productIDs, int64(productID))
	}

	// grpc request
	products, err := o.inventory.GetItems(ctx, productIDs)
	if err != nil {
		return OrderCreated{}, err
	}

	// find missing items
	found := make(map[int64]struct{})

	for _, item := range products {
		found[item.ProductID] = struct{}{}
	}

	missing := make([]int64, 0)

	for _, id := range productIDs {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}

	if len(missing) != 0 {
		return OrderCreated{}, fmt.Errorf("not found product ids: %v", missing)
	}

	issues := make([]string, 0)

	for _, inv := range products {
		for k, v := range aggregatedItems {
			if inv.ProductID == int64(k) {
				if inv.Quantity < v.Quantity {
					slog.Info("found an item but don't have enough quantity",
						"product_id", inv.ProductID,
						"quantity", inv.Quantity,
						"requested_quantity", v.Quantity,
					)

					msg := fmt.Sprintf("dont have enough quantity product id %d inv quantity %d", k, inv.Quantity)
					issues = append(issues, msg)
				} else {
					slog.Info("found an item",
						"product_id", inv.ProductID,
						"quantity", inv.Quantity,
						"unit_price", inv.UnitPrice,
					)
				}

				break
			}
		}
	}

	if len(issues) != 0 {
		var builder strings.Builder

		for _, issue := range issues {
			builder.WriteString(issue)
		}

		return OrderCreated{}, errors.New(builder.String())
	}

	order, err := o.repo.CreateOrder(products)
	if err != nil {
		return OrderCreated{}, err
	}

	return order, nil
}

// FIX: need to do grpc request to inventory to get unti price from inventory db cause we need to get total amount and etc.
//
//	{
//	  "id": 9,
//	  "status": "received",
//	  "total_amount": 82,
//	  "created_at": "2026-09-26T21:13:27.533172Z",
//	  "updated_at": "2026-09-26T21:13:27.533172Z",
//	  "items": [
//	    {
//	      "id": 7,
//	      "product_id": 0,
//	      "quantity": 3,
//	      "unit_price": 23
//	    },
//	    {
//	      "id": 8,
//	      "product_id": 1,
//	      "quantity": 1,
//	      "unit_price": 13
//	    }
//	  ]
//	}
//
// product id is not correct and id too?
// total amount isnt correct
func (o *orderService) GetOrder(ctx context.Context, orderID int) (Order, error) {
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
