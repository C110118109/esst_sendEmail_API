package user

import (
	"esst_sendEmail/internal/v1/service/user"
	model "esst_sendEmail/internal/v1/structure/users"
	"gorm.io/gorm"
)

type Resolver interface {
	Create(trx *gorm.DB, input *model.Created) interface{}
	Login(input *model.Login) interface{}
	GetByID(input *model.Field) interface{}
	List(input *model.Fields) interface{}
	Update(input *model.Updated) interface{}
	Delete(input *model.Field) interface{}
}

type resolver struct {
	UserService user.Service
}

func New(db *gorm.DB) Resolver {
	return &resolver{
		UserService: user.New(db),
	}
}
