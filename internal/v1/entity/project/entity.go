package project

import (
	model "esst_sendEmail/internal/v1/structure/projects"

	"gorm.io/gorm"
)

type Entity interface {
	WithTrx(tx *gorm.DB) Entity
	Create(input *model.Table) (err error)
	List(input *model.Fields) (int64, []*model.Table, error)
	GetByID(input *model.Field) (*model.Table, error)
	Update(input *model.Table) (err error)
	Delete(input *model.Field) (err error)
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
