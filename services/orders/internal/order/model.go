package order

import "time"

type Order struct {
	ID          int
	Status      string
	TotalAmount int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem
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

type OrderCreated struct {
	OrderID int
	Status  string
}

type Data struct {
	ProductID int
	Quantity  int
	UnitPrice int
}

type InventoryItem struct {
	ProductID int64
	Quantity  int
	UnitPrice int64
}

type DecreaseInventoryItem struct {
	ProductID int64
	Quantity  int
}
