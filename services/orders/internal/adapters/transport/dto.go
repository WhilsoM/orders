package transport

import "fmt"

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CreateOrderResponse struct {
	Status string `json:"status"`
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
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
