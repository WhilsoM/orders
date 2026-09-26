package postgresql

import (
	"backend/services/orders/internal/order"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

type orderRepository struct {
	inventory []order.Data

	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		db: db,
		inventory: []order.Data{
			{
				ProductID: 22,
				Quantity:  2,
				UnitPrice: 21,
			},
			{
				ProductID: 23,
				Quantity:  1,
				UnitPrice: 22,
			},
			{
				ProductID: 24,
				Quantity:  4,
				UnitPrice: 23,
			},
		},
	}
}

func (o *orderRepository) CreateOrder(input order.CreateOrderInput) (order.OrderCreated, error) {
	issues := make([]string, 0)

	// product_id: quantity
	aggregatedItems := make(map[int]order.Data, 0)

	slog.Info("inventory", "items", o.inventory)

	for _, val := range input.Items {
		item := aggregatedItems[val.ProductID]
		item.Quantity += val.Quantity

		aggregatedItems[val.ProductID] = item
	}

	for productID, item := range aggregatedItems {
		found := false

		for _, inv := range o.inventory {
			if inv.ProductID == productID {
				found = true

				if inv.Quantity < item.Quantity {
					slog.Info("found an item but don't have enough quantity",
						"product_id", inv.ProductID,
						"quantity", inv.Quantity,
						"requested_quantity", item.Quantity,
					)

					msg := fmt.Sprintf("don't have enough quantity product id %d inv quantity %d", productID, inv.Quantity)
					issues = append(issues, msg)
				} else {
					slog.Info("found an item",
						"product_id", inv.ProductID,
						"quantity", inv.Quantity,
						"unit_price", inv.UnitPrice,
					)

					item.UnitPrice = inv.UnitPrice
					aggregatedItems[productID] = item
				}

				break
			}
		}

		if !found {
			msg := fmt.Sprintf("not found product id %d", productID)
			issues = append(issues, msg)
		}
	}

	if len(issues) == 0 {
		orderModel := OrderModel{
			Status: "received",
		}

		if err := o.db.Create(&orderModel).Error; err != nil {
			return order.OrderCreated{}, fmt.Errorf("create order: %w", err)
		}

		items := make([]OrderItemModel, 0, len(aggregatedItems))

		for productID, item := range aggregatedItems {
			items = append(items, OrderItemModel{
				OrderID:   orderModel.ID,
				ProductID: int64(productID),
				Quantity:  item.Quantity,
				UnitPrice: float64(item.UnitPrice),
			})
		}

		if err := o.db.Create(&items).Error; err != nil {
			return order.OrderCreated{}, fmt.Errorf("create order items: %w", err)
		}

		return order.OrderCreated{
			OrderID: int(orderModel.ID),
			Status:  orderModel.Status,
		}, nil
	}

	var builder strings.Builder

	for _, issue := range issues {
		builder.WriteString(issue)
	}

	return order.OrderCreated{}, errors.New(builder.String())
}

func (o *orderRepository) GetOrder(orderID int) (order.Order, error) {
	var orderModel OrderModel

	if err := o.db.Where("id = ?", orderID).First(&orderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order.Order{}, errors.New("not found an order")
		}

		return order.Order{}, err
	}

	var orderItems []OrderItemModel

	if err := o.db.Where("order_id = ?", orderModel.ID).Find(&orderItems).Error; err != nil {
		return order.Order{}, err
	}

	result := order.Order{
		ID:        int(orderModel.ID),
		Status:    orderModel.Status,
		CreatedAt: orderModel.CreatedAt,
		UpdatedAt: orderModel.UpdatedAt,
	}

	for _, item := range orderItems {
		result.Items = append(result.Items, order.OrderItem{
			ID:        int(item.ID),
			OrderID:   int(item.OrderID),
			ProductID: int(item.ProductID),
			Quantity:  item.Quantity,
			UnitPrice: float32(item.UnitPrice),
		})
	}

	return result, nil
}
