package stock_equipment

import (
	"esst_sendEmail/internal/v1/entity/stock_equipment"
	model "esst_sendEmail/internal/v1/structure/stock_equipments"

	"gorm.io/gorm"
)

type Service interface {
	WithTrx(tx *gorm.DB) Service
	Create(input *model.Created) (*model.Base, error)
	CreateBatch(input *model.BatchCreated) error
	List(input *model.Fields) (int64, []*model.Base, error)
	ListByStockID(stockID string) ([]*model.Base, error)
	GetByID(input *model.Field) (*model.Base, error)
	Update(input *model.Updated) error
	Delete(input *model.Updated) error
}

type service struct {
	Entity stock_equipment.Entity
}

func New(db *gorm.DB) Service {
	return &service{
		Entity: stock_equipment.New(db),
	}
}

func (s *service) WithTrx(tx *gorm.DB) Service {
	return &service{
		Entity: s.Entity.WithTrx(tx),
	}
}
