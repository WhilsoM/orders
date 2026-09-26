package postgresql

import (
	"backend/services/orders/internal/order"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (o *orderRepository) CreateOrder(input []order.InventoryItem) (order.OrderCreated, error) {
	orderModel := OrderModel{
		Status: "received",
	}

	if err := o.db.Create(&orderModel).Error; err != nil {
		return order.OrderCreated{}, fmt.Errorf("create order: %w", err)
	}

	items := make([]OrderItemModel, 0, len(input))

	for productID, item := range input {
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
