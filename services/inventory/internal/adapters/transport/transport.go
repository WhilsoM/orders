package transport

type InventoryService interface {
}

type grpcServer struct {
	// inventorypb.UnimplementedInventoryServiceServer
	service InventoryService
}
