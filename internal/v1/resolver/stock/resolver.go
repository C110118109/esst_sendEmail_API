package stock

import (
	"esst_sendEmail/internal/v1/service/stock"
	"esst_sendEmail/internal/v1/service/stock_equipment"
	model "esst_sendEmail/internal/v1/structure/stocks"

	"gorm.io/gorm"
)

type Resolver interface {
	Create(trx *gorm.DB, input *model.Created) interface{}
	List(input *model.Fields) interface{}
	GetByID(input *model.Field) interface{}
	Update(input *model.Updated) interface{}
	Delete(input *model.Updated) interface{}
}

type resolver struct {
	StockService          stock.Service
	StockEquipmentService stock_equipment.Service
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		StockService:          stock.New(db),
		StockEquipmentService: stock_equipment.New(db),
	}
}
