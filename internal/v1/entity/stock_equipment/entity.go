package stock_equipment

import (
	model "esst_sendEmail/internal/v1/structure/stock_equipments"

	"gorm.io/gorm"
)

type Entity interface {
	WithTrx(tx *gorm.DB) Entity
	Create(input *model.Table) (err error)
	CreateBatch(input *model.BatchCreated) (err error)
	List(input *model.Fields) (int64, []*model.Table, error)
	ListByStockID(stockID string) ([]*model.Table, error)
	GetByID(input *model.Field) (*model.Table, error)
	Update(input *model.Table) (err error)
	Delete(input *model.Field) (err error)
	DeleteByStockID(stockID string) (err error)
}

type entity struct {
	db *gorm.DB
}

func New(db *gorm.DB) Entity {
	return &entity{db: db}
}

func (e *entity) WithTrx(tx *gorm.DB) Entity {
	return &entity{db: tx}
}
