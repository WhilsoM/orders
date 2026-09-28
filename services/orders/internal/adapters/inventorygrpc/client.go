package inventorygrpc

import (
	inventorypb "backend/pkg/gen/inventory"
	"backend/services/orders/internal/order"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client inventorypb.InventoryServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: inventorypb.NewInventoryServiceClient(conn),
	}, nil
}

func (c *Client) GetItems(
	ctx context.Context,
	productIDs []int64,
) ([]order.InventoryItem, error) {
	resp, err := c.client.GetItems(
		ctx,
		&inventorypb.GetItemsRequest{
			ProductIds: productIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	items := make([]order.InventoryItem, 0, len(resp.Items))

	for _, item := range resp.Items {
		items = append(items, order.InventoryItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			UnitPrice: item.UnitPrice,
		})
	}

	return items, nil
}

func (c *Client) DecreaseItems(
	ctx context.Context,
	items []order.DecreaseInventoryItem,
) error {
	pbItems := make([]*inventorypb.DecreaseItem, 0, len(items))

	for _, item := range items {
		pbItems = append(pbItems, &inventorypb.DecreaseItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	_, err := c.client.DecreaseItems(
		ctx,
		&inventorypb.DecreaseItemsRequest{
			Items: pbItems,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
