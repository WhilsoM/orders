package postgresql

import (
	"backend/services/inventory/internal/inventory"
	"fmt"

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

func (r *inventoryRepository) DecreaseItems(items []inventory.DecreaseItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.
				Model(&inventory.ItemModel{}).
				Where(
					"product_id = ? AND quantity >= ?",
					item.ProductID,
					item.Quantity,
				).
				Update(
					"quantity",
					gorm.Expr("quantity - ?", item.Quantity),
				)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf(
					"cannot decrease product %d by %d",
					item.ProductID,
					item.Quantity,
				)
			}
		}

		return nil
	})
}
