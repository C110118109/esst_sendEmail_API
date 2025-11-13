package stock

import (
	model "esst_sendEmail/internal/v1/structure/stocks"
)

func (e *entity) Create(input *model.Table) error {
	return e.db.Create(input).Error
}

func (e *entity) List(input *model.Fields) (int64, []*model.Table, error) {
	var total int64
	var records []*model.Table

	db := e.db.Model(&model.Table{})

	// 篩選條件
	if input.StockName != nil {
		db = db.Where("stock_name LIKE ?", "%"+*input.StockName+"%")
	}
	if input.ContactName != nil {
		db = db.Where("contact_name LIKE ?", "%"+*input.ContactName+"%")
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
	db := e.db.Model(&model.Table{}).Where("stock_id = ?", input.StockID)
	err = db.First(&output).Error
	return output, err
}

func (e *entity) Update(input *model.Table) (err error) {
	err = e.db.Model(&model.Table{}).Where("stock_id = ?", input.StockID).Updates(input).Error
	return err
}

func (e *entity) Delete(input *model.Field) error {
	return e.db.Where("stock_id = ?", input.StockID).Delete(&model.Table{}).Error
}
