package equipment

import (
	"esst_sendEmail/internal/v1/service/equipment"
	"esst_sendEmail/internal/v1/service/project"
	model "esst_sendEmail/internal/v1/structure/equipments"
	projectModel "esst_sendEmail/internal/v1/structure/projects"

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
	ProjectService   project.Service
}

// Field 用於查詢專案
type Field struct {
	ProjectID string
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		EquipmentService: equipment.New(db),
		ProjectService:   project.New(db),
	}
}

// GetByID 的輔助方法
func (f *Field) ToProjectField() *projectModel.Field {
	return &projectModel.Field{
		ProjectID: f.ProjectID,
	}
}
