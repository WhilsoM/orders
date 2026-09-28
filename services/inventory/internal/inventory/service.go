package inventory

import "fmt"

type InventoryRepository interface {
	GetItems(productIDs []int64) ([]ItemModel, error)
	DecreaseItems(items []DecreaseItem) error
}

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *inventoryService {
	return &inventoryService{
		repo: repo,
	}
}

func (s *inventoryService) GetItems(productIDs []int64) ([]Item, error) {
	if len(productIDs) == 0 {
		return []Item{}, nil
	}

	products, err := s.repo.GetItems(productIDs)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(products))

	for i, v := range products {
		items[i] = Item{
			ProductID: v.ProductID,
			Quantity:  v.Quantity,
			UnitPrice: v.UnitPrice,
		}
	}

	return items, nil
}

func (s *inventoryService) DecreaseItems(items []DecreaseItem) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			return fmt.Errorf(
				"invalid quantity for product %d",
				item.ProductID,
			)
		}
	}

	return s.repo.DecreaseItems(items)
}
