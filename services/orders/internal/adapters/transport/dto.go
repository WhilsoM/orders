package transport

import (
	"fmt"
	"time"
)

type Order struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	Order
}

type CreateOrderResponse struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

func (r CreateOrderRequest) Validate() error {
	if len(r.Items) == 0 {
		return fmt.Errorf("items must not be empty")
	}

	for i, item := range r.Items {
		if err := item.Validate(i); err != nil {
			return err
		}
	}

	return nil
}

func (i CreateOrderItemRequest) Validate(index int) error {
	if i.ProductID <= 0 {
		return fmt.Errorf(
			"items %d product_id must be greater than zero",
			index,
		)
	}

	if i.Quantity <= 0 {
		return fmt.Errorf(
			"items %d quantity must be greater than zero",
			index,
		)
	}

	return nil
}

type GetOrderByIDResponse struct {
	ID          int                `json:"id"`
	Status      string             `json:"status"`
	TotalAmount int                `json:"total_amount"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Items       []GetOrderByIDItem `json:"items"`
}

type GetOrderByIDItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float32 `json:"unit_price"`
}
