package postgresql

type ItemModel struct {
	ProductID int64
	Quantity  int
	UnitPrice int
}

func (i ItemModel) TableName() string {
	return "inventory"
}
