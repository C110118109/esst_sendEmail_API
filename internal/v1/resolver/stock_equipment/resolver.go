package stock_equipment

import (
	"esst_sendEmail/internal/v1/service/stock"
	"esst_sendEmail/internal/v1/service/stock_equipment"
	model "esst_sendEmail/internal/v1/structure/stock_equipments"
	stockModel "esst_sendEmail/internal/v1/structure/stocks"

	"gorm.io/gorm"
)

type Resolver interface {
	Create(trx *gorm.DB, input *model.Created) interface{}
	CreateBatch(trx *gorm.DB, input *model.BatchCreated) interface{}
	List(input *model.Fields) interface{}
	ListByStockID(stockID string) interface{}
	GetByID(input *model.Field) interface{}
	Update(input *model.Updated) interface{}
	Delete(input *model.Updated) interface{}
}

type resolver struct {
	StockEquipmentService stock_equipment.Service
	StockService          stock.Service
}

// Field 用於查詢現貨
type Field struct {
	StockID string
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		StockEquipmentService: stock_equipment.New(db),
		StockService:          stock.New(db),
	}
}

// GetByID 的輔助方法
func (f *Field) ToStockField() *stockModel.Field {
	return &stockModel.Field{
		StockID: f.StockID,
	}
}
