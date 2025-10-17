package project

import (
	"esst_sendEmail/internal/v1/entity/project"
	model "esst_sendEmail/internal/v1/structure/projects"

	"gorm.io/gorm"
)

type Service interface {
	WithTrx(tx *gorm.DB) Service
	Create(input *model.Created) (*model.Base, error)
	List(input *model.Fields) (int64, []*model.Base, error)
	GetByID(input *model.Field) (*model.Base, error)
	Update(input *model.Updated) error
	Delete(input *model.Updated) error
}

type service struct {
	Entity project.Entity
}

func New(db *gorm.DB) Service {
	return &service{
		Entity: project.New(db),
	}
}

func (s *service) WithTrx(tx *gorm.DB) Service {
	return &service{
		Entity: s.Entity.WithTrx(tx),
	}
}
