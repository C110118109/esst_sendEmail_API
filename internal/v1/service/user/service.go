package user

import (
	"esst_sendEmail/internal/v1/entity/user"
	model "esst_sendEmail/internal/v1/structure/users"
	"gorm.io/gorm"
)

type Service interface {
	WithTrx(tx *gorm.DB) Service
	Create(input *model.Created) (*model.Base, error)
	Authenticate(username, password string) (*model.Base, error)
	GetByID(input *model.Field) (*model.Base, error)
	List(input *model.Fields) (int64, []*model.Base, error)
	Update(input *model.Updated) error
	Delete(input *model.Field) error
}

type service struct {
	Entity user.Entity
}

func New(db *gorm.DB) Service {
	return &service{
		Entity: user.New(db),
	}
}

func (s *service) WithTrx(tx *gorm.DB) Service {
	return &service{
		Entity: s.Entity.WithTrx(tx),
	}
}
