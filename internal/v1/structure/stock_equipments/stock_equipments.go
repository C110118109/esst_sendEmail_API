package stock_equipments

import (
	model "esst_sendEmail/internal/v1/structure"
	"esst_sendEmail/internal/v1/structure/stocks"
	"time"
)

// StockEquipment struct is a row record of the stock_equipments table
// Table struct is database table struct
type Table struct {
	// 現貨設備編號
	StockEquipmentID string `gorm:"primaryKey;uuid_generate_v4();column:seq_id;type:uuid;" json:"seq_id,omitempty"`
	// 現貨編號
	StockID string `gorm:"column:stock_id;type:uuid;not null" json:"stock_id"`
	// stock data
	Stock stocks.Table `gorm:"foreignKey:StockID;references:StockID" json:"-"`
	// 料號
	PartNumber string `gorm:"column:part_number;type:TEXT;" json:"part_number,omitempty"`
	// 數量
	Quantity int64 `gorm:"column:quantity;type:integer;" json:"quantity,omitempty"`
	// 說明
	Description string `gorm:"column:description;type:TEXT;" json:"description,omitempty"`
	// 創建時間
	CreatedTime time.Time `gorm:"column:created_time;type:TIMESTAMP;" json:"created_time"`
}

// Base struct is corresponding to table structure file
type Base struct {
	// 現貨設備編號
	StockEquipmentID string `json:"seq_id,omitempty"`
	// 現貨編號
	StockID string `json:"stock_id"`
	// 料號
	PartNumber string `json:"part_number,omitempty"`
	// 數量
	Quantity int64 `json:"quantity,omitempty"`
	// 說明
	Description string `json:"description,omitempty"`
	// 創建時間
	CreatedTime time.Time `json:"created_time"`
}

// Single return structure file
type Single struct {
	// 現貨設備編號
	StockEquipmentID string `json:"seq_id,omitempty"`
	// 現貨編號
	StockID string `json:"stock_id"`
	// 料號
	PartNumber string `json:"part_number,omitempty"`
	// 數量
	Quantity int64 `json:"quantity,omitempty"`
	// 說明
	Description string `json:"description,omitempty"`
	// 創建時間
	CreatedTime time.Time `json:"created_time"`
}

// Created struct is used to create
type Created struct {
	// 現貨編號
	StockID string `json:"stock_id" binding:"required,uuid4" form:"stock_id"`
	// 料號
	PartNumber string `json:"part_number" binding:"required" validate:"required"`
	// 數量
	Quantity int64 `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	// 說明
	Description string `json:"description,omitempty"`
}

// BatchCreated struct is used to create multiple stock equipments
type BatchCreated struct {
	// 現貨編號
	StockID string `json:"stock_id" binding:"required,uuid4" validate:"required,uuid4"`
	// 設備列表
	Equipments []BatchEquipment `json:"equipments" binding:"required,min=1,dive" validate:"required,min=1,dive"`
}

// BatchEquipment struct for batch creation
type BatchEquipment struct {
	// 料號
	PartNumber string `json:"part_number" binding:"required" validate:"required"`
	// 數量
	Quantity int64 `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	// 說明
	Description string `json:"description,omitempty"`
}

// Field is structure file for search
type Field struct {
	// 現貨設備編號
	StockEquipmentID string `json:"seq_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	// 現貨編號
	StockID string `json:"stock_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	// 料號
	PartNumber string `json:"part_number,omitempty"`
	// 數量
	Quantity int64 `json:"quantity,omitempty"`
	// 說明
	Description string `json:"description,omitempty"`
}

// Fields is the searched structure file (including pagination)
type Fields struct {
	Field
	model.InPage
}

// List is multiple return structure files
type List struct {
	StockEquipments []*struct {
		// 現貨設備編號
		StockEquipmentID string `json:"seq_id,omitempty"`
		// 現貨編號
		StockID string `json:"stock_id"`
		// 料號
		PartNumber string `json:"part_number,omitempty"`
		// 數量
		Quantity int64 `json:"quantity,omitempty"`
		// 說明
		Description string `json:"description,omitempty"`
		// 創建時間
		CreatedTime time.Time `json:"created_time"`
	} `json:"stock_equipments"`
	model.OutPage
}

// Updated struct is used to update
type Updated struct {
	// 現貨設備編號
	StockEquipmentID string `json:"seq_id" binding:"required,uuid4" validate:"required,uuid4"`
	// 現貨編號
	StockID string `json:"stock_id" binding:"required,uuid4" validate:"required,uuid4"`
	// 料號
	PartNumber string `json:"part_number" binding:"required" validate:"required"`
	// 數量
	Quantity int64 `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	// 說明
	Description string `json:"description,omitempty"`
}

// TableName sets the insert table name for this struct type
func (a *Table) TableName() string {
	return "stock_equipments"
}
