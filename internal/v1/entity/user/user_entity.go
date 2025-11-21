package user

import (
	"errors"
	model "esst_sendEmail/internal/v1/structure/users"
)

func (e *entity) Create(input *model.Table) error {
	return e.db.Create(input).Error
}

func (e *entity) GetByUsername(input *model.Field) (*model.Table, error) {
	var output model.Table
	err := e.db.Where("username = ?", *input.Username).First(&output).Error
	return &output, err
}

func (e *entity) GetByID(input *model.Field) (*model.Table, error) {
	var output model.Table
	err := e.db.Where("id = ?", *input.ID).First(&output).Error
	return &output, err
}

func (e *entity) List(input *model.Fields) (int64, []*model.Table, error) {
	var total int64
	var records []*model.Table

	db := e.db.Model(&model.Table{})

	if input.Username != nil {
		db = db.Where("username LIKE ?", "%"+*input.Username+"%")
	}

	if input.Email != nil {
		db = db.Where("email LIKE ?", "%"+*input.Email+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Order("created_at DESC").
		Offset(int((input.Page - 1) * input.Limit)).
		Limit(int(input.Limit)).
		Find(&records).Error

	return total, records, err
}

func (e *entity) Update(input *model.Table) error {
	return e.db.Model(&model.Table{}).Where("id = ?", input.ID).Updates(input).Error
}

func (e *entity) Delete(input *model.Field) error {
	//後端保護:先查詢使用者資料
	var user model.Table
	err := e.db.Where("id = ?", *input.ID).First(&user).Error
	if err != nil {
		return err
	}

	// 禁止刪除預設管理員帳號
	if user.Username == "admin" {
		return errors.New("無法刪除預設管理員帳號")
	}

	// 執行刪除
	return e.db.Where("id = ?", *input.ID).Delete(&model.Table{}).Error
}
