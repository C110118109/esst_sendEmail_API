package stock_equipment

import (
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/stock_equipments"
	"fmt"
	"time"
)

// 批次建立現貨設備
func (e *entity) CreateBatch(input *model.BatchCreated) error {
	// 驗證現貨 ID
	if input.StockID == "" {
		log.Error("StockID is empty")
		return fmt.Errorf("stock ID is required")
	}

	// 將每個設備資料轉換為 Table 結構
	var equipmentTables []*model.Table
	for _, eq := range input.Equipments {
		equipment := &model.Table{
			StockEquipmentID: util.GenerateUUID(),
			StockID:          input.StockID,
			PartNumber:       eq.PartNumber,
			Quantity:         eq.Quantity,
			Description:      eq.Description,
			CreatedTime:      time.Now(),
		}
		equipmentTables = append(equipmentTables, equipment)
	}

	// 使用批次建立
	if len(equipmentTables) > 0 {
		if err := e.db.Create(&equipmentTables).Error; err != nil {
			log.Error("Failed to insert stock equipments:", err)
			return err
		}
	}

	return nil
}

// 根據現貨 ID 獲取設備列表
func (e *entity) ListByStockID(stockID string) ([]*model.Table, error) {
	var equipments []*model.Table
	err := e.db.Where("stock_id = ?", stockID).Find(&equipments).Error
	if err != nil {
		log.Error("Failed to query stock equipments by stock ID:", err)
		return nil, err
	}
	return equipments, err
}

// 原有的函數保持不變
func (e *entity) Create(input *model.Table) error {
	return e.db.Create(input).Error
}

func (e *entity) List(input *model.Fields) (int64, []*model.Table, error) {
	var total int64
	var records []*model.Table

	db := e.db.Model(&model.Table{})

	if input.PartNumber != "" {
		db = db.Where("part_number LIKE ?", "%"+input.PartNumber+"%")
	}

	if input.StockID != "" {
		db = db.Where("stock_id = ?", input.StockID)
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
	db := e.db.Model(&model.Table{}).Where("seq_id = ?", input.StockEquipmentID)
	err = db.First(&output).Error
	return output, err
}

func (e *entity) Update(input *model.Table) (err error) {
	err = e.db.Model(&model.Table{}).Save(&input).Error
	return err
}

func (e *entity) Delete(input *model.Field) error {
	return e.db.Where("seq_id = ?", input.StockEquipmentID).Delete(&model.Table{}).Error
}

// 根據現貨 ID 刪除所有相關設備
func (e *entity) DeleteByStockID(stockID string) error {
	return e.db.Where("stock_id = ?", stockID).Delete(&model.Table{}).Error
}
