package project

import (
	model "esst_sendEmail/internal/v1/structure/projects"
)

func (e *entity) Create(input *model.Table) error {
	return e.db.Create(input).Error
}

func (e *entity) List(input *model.Fields) (int64, []*model.Table, error) {
	var total int64
	var records []*model.Table

	db := e.db.Model(&model.Table{})

	if input.ProjectName != nil {
		db = db.Where("p_name LIKE ?", "%"+*input.ProjectName+"%")
	}
	if input.ContactName != nil {
		db = db.Where("contact_name LIKE ?", "%"+*input.ContactName+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Offset(int((input.Page - 1) * input.Limit)).
		Limit(int(input.Limit)).
		Find(&records).Error

	return total, records, err
}

func (e *entity) GetByID(input *model.Field) (output *model.Table, err error) {
	db := e.db.Model(&model.Table{}).Where("p_id = ?", input.ProjectID)

	err = db.First(&output).Error

	return output, err
}

func (e *entity) Update(input *model.Table) (err error) {
	err = e.db.Model(&model.Table{}).Save(&input).Error

	return err
}

func (e *entity) Delete(input *model.Field) error {
	return e.db.Where("p_id = ?", input.ProjectID).Delete(&model.Table{}).Error
}
