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

	// 第一階段欄位篩選
	if input.ProjectName != nil {
		db = db.Where("p_name LIKE ?", "%"+*input.ProjectName+"%")
	}
	if input.ContactName != nil {
		db = db.Where("contact_name LIKE ?", "%"+*input.ContactName+"%")
	}
	if input.ContactPhone != nil {
		db = db.Where("contact_phone LIKE ?", "%"+*input.ContactPhone+"%")
	}
	if input.ContactEmail != nil {
		db = db.Where("contact_email LIKE ?", "%"+*input.ContactEmail+"%")
	}
	if input.Owner != nil {
		db = db.Where("owner LIKE ?", "%"+*input.Owner+"%")
	}

	// 專案狀態篩選
	if input.Status != nil {
		db = db.Where("status = ?", *input.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	// 按更新時間降序排列，若無更新時間則按創建時間降序
	err = db.Order("COALESCE(updated_time, created_time) DESC").
		Offset(int((input.Page - 1) * input.Limit)).
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
	// 使用 Save 會更新所有欄位（包括零值）
	// 如果只想更新非零值欄位，可以使用 Updates
	err = e.db.Model(&model.Table{}).Where("p_id = ?", input.ProjectID).Updates(input).Error

	return err
}

func (e *entity) Delete(input *model.Field) error {
	return e.db.Where("p_id = ?", input.ProjectID).Delete(&model.Table{}).Error
}
