package equipments

import (
	model "esst_sendEmail/internal/v1/structure"
	"esst_sendEmail/internal/v1/structure/projects"
	"time"
)

// Equipment struct is a row record of the companies table in the invoice database
// Table struct is database table struct
type Table struct {
	// 設備編號
	EquipmentID string `gorm:"primaryKey;uuid_generate_v4();column:eq_id;type:uuid;" json:"eq_id,omitempty"`
	// 專案編號
	ProjectID string `gorm:"column:p_id;type:uuid;not null" json:"p_id"`
	// project data
	Project projects.Table `gorm:"foreignKey:ProjectID;references:ProjectID" json:"-"`
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
	// 設備編號
	EquipmentID string `json:"eq_id,omitempty"`
	// 專案編號
	ProjectID string `json:"p_id"`
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
	// 設備編號
	EquipmentID string `json:"eq_id,omitempty"`
	// 專案編號
	ProjectID string `json:"p_id"`
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
	// 專案編號
	ProjectID string `json:"p_id" binding:"required,uuid4" form:"p_id"`
	// 料號
	PartNumber string `json:"part_number" binding:"required" validate:"required"`
	// 數量
	Quantity int64 `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	// 說明
	Description string `json:"description,omitempty"`
}

// BatchCreated struct is used to create multiple equipments
type BatchCreated struct {
	// 專案編號
	ProjectID string `json:"p_id" binding:"required,uuid4" validate:"required,uuid4"`
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
	// 設備編號
	EquipmentID string `json:"eq_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	// 專案編號
	ProjectID string `json:"p_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
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
	Equipments []*struct {
		// 設備編號
		EquipmentID string `json:"eq_id,omitempty"`
		// 專案編號
		ProjectID string `json:"p_id"`
		// 料號
		PartNumber string `json:"part_number,omitempty"`
		// 數量
		Quantity int64 `json:"quantity,omitempty"`
		// 說明
		Description string `json:"description,omitempty"`
		// 創建時間
		CreatedTime time.Time `json:"created_time"`
	} `json:"equipments"`
	model.OutPage
}

// Updated struct is used to update
type Updated struct {
	// 設備編號
	EquipmentID string `json:"eq_id" binding:"required,uuid4" validate:"required,uuid4"`
	// 專案編號
	ProjectID string `json:"p_id" binding:"required,uuid4" validate:"required,uuid4"`
	// 料號
	PartNumber string `json:"part_number" binding:"required" validate:"required"`
	// 數量
	Quantity int64 `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	// 說明
	Description string `json:"description,omitempty"`
}

// TableName sets the insert table name for this struct type
func (a *Table) TableName() string {
	return "equipments"
}
