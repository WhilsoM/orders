package postgresql

import "time"

type OrderModel struct {
	ID        int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (OrderModel) TableName() string {
	return "orders"
}

type OrderItemModel struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	UnitPrice float64
}

func (OrderItemModel) TableName() string {
	return "order_items"
}

type GetOrderByIDResponse struct {
	OrderID int64
	Items   []OrderItemModel
}
