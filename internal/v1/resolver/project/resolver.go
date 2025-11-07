package project

import (
	"esst_sendEmail/internal/v1/service/equipment"
	"esst_sendEmail/internal/v1/service/project"
	model "esst_sendEmail/internal/v1/structure/projects"

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
	ProjectService   project.Service
	EquipmentService equipment.Service
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		ProjectService:   project.New(db),
		EquipmentService: equipment.New(db),
	}
}
