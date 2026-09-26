package postgresql

import (
	"backend/services/inventory/internal/inventory"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) GetItems(productIDs []int64) ([]inventory.ItemModel, error) {
	var items []inventory.ItemModel

	if err := r.db.
		Where("product_id IN ?", productIDs).
		Find(&items).
		Error; err != nil {
		return nil, err
	}

	return items, nil
}
