package transport

import (
	inventorypb "backend/pkg/gen/inventory"
	"backend/services/inventory/internal/inventory"
	"context"
)

type InventoryService interface {
	GetItems(productIDs []int64) ([]inventory.Item, error)
	DecreaseItems(items []inventory.DecreaseItem) error
}

type grpcServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	service InventoryService
}

func NewGRPCServer(service InventoryService) *grpcServer {
	return &grpcServer{
		service: service,
	}
}

func (s *grpcServer) GetItems(
	ctx context.Context,
	req *inventorypb.GetItemsRequest,
) (*inventorypb.GetItemsResponse, error) {
	items, err := s.service.GetItems(req.ProductIds)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*inventorypb.Item, 0, len(items))

	for _, item := range items {
		pbItems = append(pbItems, &inventorypb.Item{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			UnitPrice: int64(item.UnitPrice),
		})
	}

	return &inventorypb.GetItemsResponse{
		Items: pbItems,
	}, nil
}

func (s *grpcServer) DecreaseItems(
	ctx context.Context,
	req *inventorypb.DecreaseItemsRequest,
) (*inventorypb.DecreaseItemsResponse, error) {
	items := make([]inventory.DecreaseItem, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, inventory.DecreaseItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		})
	}

	if err := s.service.DecreaseItems(items); err != nil {
		return nil, err
	}

	return &inventorypb.DecreaseItemsResponse{}, nil
}
