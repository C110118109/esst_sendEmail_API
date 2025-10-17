package equipment

import (
	"esst_sendEmail/internal/v1/service/equipment"
	model "esst_sendEmail/internal/v1/structure/equipments"

	"gorm.io/gorm"
)

type Resolver interface {
	Create(trx *gorm.DB, input *model.Created) interface{}
	CreateBatch(trx *gorm.DB, input *model.BatchCreated) interface{}
	List(input *model.Fields) interface{}
	ListByProjectID(projectID string) interface{}
	GetByID(input *model.Field) interface{}
	Update(input *model.Updated) interface{}
	Delete(input *model.Updated) interface{}
}

type resolver struct {
	EquipmentService equipment.Service
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		EquipmentService: equipment.New(db),
	}
}
