package postgresql

import (
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

func (r *inventoryRepository) GetItems(productIDs []int64) ([]ItemModel, error) {
	var items []ItemModel

	if err := r.db.
		Where("product_id IN ?", productIDs).
		Find(&items).
		Error; err != nil {
		return nil, err
	}

	return items, nil
}
