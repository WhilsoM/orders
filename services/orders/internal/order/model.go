package order

import "time"

type Order struct {
	ID          int
	Status      string
	TotalAmount int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	UnitPrice float32
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID int
	Quantity  int
}

type Data struct {
	ProductID int
	Quantity  int
}
