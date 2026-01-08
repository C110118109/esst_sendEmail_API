package projects

import (
	model "esst_sendEmail/internal/v1/structure"
	"time"
)

// Poject struct is a row record of the companies table in the invoice database
// Table struct is database table struct
type Table struct {
	// 專案編號
	ProjectID string `gorm:"primaryKey;uuid_generate_v4();column:p_id;type:uuid;" json:"p_id,omitempty"`
	// 專案名稱
	ProjectName string `gorm:"column:p_name;type:TEXT;" json:"p_name,omitempty"`
	// 聯絡姓名
	ContactName string `gorm:"column:contact_name;type:TEXT;" json:"contact_name,omitempty"`
	// 聯絡電話
	ContactPhone string `gorm:"column:contact_phone;type:TEXT;" json:"contact_phone,omitempty"`
	// 聯絡信箱
	ContactEmail string `gorm:"column:contact_email;type:TEXT;" json:"contact_email,omitempty"`
	// 雙欣負責人
	Owner string `gorm:"column:owner;type:TEXT;" json:"owner,omitempty"`
	// 備註
	Remark string `gorm:"column:remark;type:TEXT;" json:"remark,omitempty"`

	// 第二階段欄位
	// 預計交貨期
	ExpectedDeliveryPeriod string `gorm:"column:expected_delivery_period;type:TEXT;" json:"expected_delivery_period,omitempty"`
	// 預計交貨日
	ExpectedDeliveryDate *time.Time `gorm:"column:expected_delivery_date;type:DATE;" json:"expected_delivery_date,omitempty"`
	// 預計履約期
	ExpectedContractPeriod string `gorm:"column:expected_contract_period;type:TEXT;" json:"expected_contract_period,omitempty"`
	// 履約開始日
	ContractStartDate *time.Time `gorm:"column:contract_start_date;type:DATE;" json:"contract_start_date,omitempty"`
	// 履約結束日
	ContractEndDate *time.Time `gorm:"column:contract_end_date;type:DATE;" json:"contract_end_date,omitempty"`
	// 交貨地址
	DeliveryAddress string `gorm:"column:delivery_address;type:TEXT;" json:"delivery_address,omitempty"`
	// 特殊需求
	SpecialRequirements string `gorm:"column:special_requirements;type:TEXT;" json:"special_requirements,omitempty"`
	// 專案狀態 (step1: 第一階段, step2: 第二階段, completed: 已完成)
	Status string `gorm:"column:status;type:TEXT;default:'step1';" json:"status,omitempty"`

	// 創建時間
	CreatedTime time.Time `gorm:"column:created_time;type:TIMESTAMP;" json:"created_time"`
	// 更新時間
	UpdatedTime *time.Time `gorm:"column:updated_time;type:TIMESTAMP;" json:"updated_time,omitempty"`
}

// Base struct is corresponding to table structure file
type Base struct {
	// 專案編號
	ProjectID string `json:"p_id,omitempty"`
	// 專案名稱
	ProjectName string `json:"p_name,omitempty"`
	// 聯絡姓名
	ContactName string `json:"contact_name,omitempty"`
	// 聯絡電話
	ContactPhone string `json:"contact_phone,omitempty"`
	// 聯絡信箱
	ContactEmail string `json:"contact_email,omitempty"`
	// 雙欣負責人
	Owner string `json:"owner,omitempty"`
	// 備註
	Remark string `json:"remark,omitempty"`

	// 第二階段欄位
	// 預計交貨期
	ExpectedDeliveryPeriod string `json:"expected_delivery_period,omitempty"`
	// 預計交貨日
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	// 預計履約期
	ExpectedContractPeriod string `json:"expected_contract_period,omitempty"`
	// 履約開始日
	ContractStartDate *time.Time `json:"contract_start_date,omitempty"`
	// 履約結束日
	ContractEndDate *time.Time `json:"contract_end_date,omitempty"`
	// 交貨地址
	DeliveryAddress string `json:"delivery_address,omitempty"`
	// 特殊需求
	SpecialRequirements string `json:"special_requirements,omitempty"`
	// 專案狀態
	Status string `json:"status,omitempty"`

	// 創建時間
	CreatedTime time.Time `json:"created_time"`
	// 更新時間
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
}

// Single return structure file
type Single struct {
	// 專案編號
	ProjectID string `json:"p_id,omitempty"`
	// 專案名稱
	ProjectName string `json:"p_name,omitempty"`
	// 聯絡姓名
	ContactName string `json:"contact_name,omitempty"`
	// 聯絡電話
	ContactPhone string `json:"contact_phone,omitempty"`
	// 聯絡信箱
	ContactEmail string `json:"contact_email,omitempty"`
	// 雙欣負責人
	Owner string `json:"owner,omitempty"`
	// 備註
	Remark string `json:"remark,omitempty"`

	// 第二階段欄位
	// 預計交貨期
	ExpectedDeliveryPeriod string `json:"expected_delivery_period,omitempty"`
	// 預計交貨日
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	// 預計履約期
	ExpectedContractPeriod string `json:"expected_contract_period,omitempty"`
	// 履約開始日
	ContractStartDate *time.Time `json:"contract_start_date,omitempty"`
	// 履約結束日
	ContractEndDate *time.Time `json:"contract_end_date,omitempty"`
	// 交貨地址
	DeliveryAddress string `json:"delivery_address,omitempty"`
	// 特殊需求
	SpecialRequirements string `json:"special_requirements,omitempty"`
	// 專案狀態
	Status string `json:"status,omitempty"`

	// 創建時間
	CreatedTime time.Time `json:"created_time"`
	// 更新時間
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
}

// Created struct is used to create
type Created struct {
	// 專案名稱
	ProjectName string `json:"p_name" binding:"required" validate:"required"`
	// 聯絡姓名
	ContactName string `json:"contact_name" binding:"required" validate:"required"`
	// 聯絡電話
	ContactPhone string `json:"contact_phone" binding:"required" validate:"required"`
	// 聯絡信箱
	ContactEmail string `json:"contact_email" binding:"required,email" validate:"email"`
	// 雙欣負責人
	Owner string `json:"owner" binding:"required" validate:"required"`
	// 備註
	Remark string `json:"remark,omitempty"`
}

// Field is structure file for search
type Field struct {
	// 專案編號
	ProjectID string `json:"p_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	// 專案名稱
	ProjectName *string `json:"p_name,omitempty" form:"p_name"`
	// 聯絡姓名
	ContactName *string `json:"contact_name,omitempty" form:"contact_name"`
	// 聯絡電話
	ContactPhone *string `json:"contact_phone,omitempty" form:"contact_phone"`
	// 聯絡信箱
	ContactEmail *string `json:"contact_email,omitempty" form:"contact_email"`
	// 雙欣負責人
	Owner *string `json:"owner,omitempty" form:"owner"`
	// 專案狀態
	Status *string `json:"status,omitempty" form:"status"`
}

// Fields is the searched structure file (including pagination)
type Fields struct {
	Field
	model.InPage
}

// List is multiple return structure files
type List struct {
	Projects []*struct {
		// 專案編號
		ProjectID string `json:"p_id,omitempty"`
		// 專案名稱
		ProjectName string `json:"p_name,omitempty"`
		// 聯絡姓名
		ContactName string `json:"contact_name,omitempty"`
		// 聯絡電話
		ContactPhone string `json:"contact_phone,omitempty"`
		// 聯絡信箱
		ContactEmail string `json:"contact_email,omitempty"`
		// 雙欣負責人
		Owner string `json:"owner,omitempty"`
		// 備註
		Remark string `json:"remark,omitempty"`

		// 第二階段欄位
		// 預計交貨期
		ExpectedDeliveryPeriod string `json:"expected_delivery_period,omitempty"`
		// 預計交貨日
		ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
		// 預計履約期
		ExpectedContractPeriod string `json:"expected_contract_period,omitempty"`
		// 履約開始日
		ContractStartDate *time.Time `json:"contract_start_date,omitempty"`
		// 履約結束日
		ContractEndDate *time.Time `json:"contract_end_date,omitempty"`
		// 交貨地址
		DeliveryAddress string `json:"delivery_address,omitempty"`
		// 特殊需求
		SpecialRequirements string `json:"special_requirements,omitempty"`
		// 專案狀態
		Status string `json:"status,omitempty"`

		// 創建時間
		CreatedTime time.Time `json:"created_time"`
		// 更新時間
		UpdatedTime *time.Time `json:"updated_time,omitempty"`
	} `json:"projects"`
	model.OutPage
}

// Updated struct is used to update
type Updated struct {
	// 專案編號
	ProjectID string `json:"p_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	// 專案名稱
	ProjectName string `json:"p_name,omitempty"`
	// 聯絡姓名
	ContactName string `json:"contact_name,omitempty"`
	// 聯絡電話
	ContactPhone string `json:"contact_phone,omitempty"`
	// 聯絡信箱
	ContactEmail string `json:"contact_email,omitempty"`
	// 雙欣負責人
	Owner string `json:"owner,omitempty"`
	// 備註
	Remark string `json:"remark,omitempty"`

	// 第二階段欄位
	// 預計交貨期
	ExpectedDeliveryPeriod string `json:"expected_delivery_period,omitempty"`
	// 預計交貨日（接收字串格式，如 "2025-10-23"）
	ExpectedDeliveryDate string `json:"expected_delivery_date,omitempty"`
	// 預計履約期
	ExpectedContractPeriod string `json:"expected_contract_period,omitempty"`
	// 履約開始日（接收字串格式，如 "2025-10-23"）
	ContractStartDate string `json:"contract_start_date,omitempty"`
	// 履約結束日（接收字串格式，如 "2025-10-23"）
	ContractEndDate string `json:"contract_end_date,omitempty"`
	// 交貨地址
	DeliveryAddress string `json:"delivery_address,omitempty"`
	// 特殊需求
	SpecialRequirements string `json:"special_requirements,omitempty"`
	// 專案狀態
	Status string `json:"status,omitempty"`
}

// TableName sets the insert table name for this struct type
func (a *Table) TableName() string {
	return "projects"
}
