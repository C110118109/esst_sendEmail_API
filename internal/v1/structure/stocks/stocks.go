package stocks

import (
	model "esst_sendEmail/internal/v1/structure"
	"time"
)

// Stock struct is a row record of the stocks table
// Table struct is database table struct
type Table struct {
	// 現貨編號
	StockID string `gorm:"primaryKey;uuid_generate_v4();column:stock_id;type:uuid;" json:"stock_id,omitempty"`
	
	// 基本資訊
	StockName    string `gorm:"column:stock_name;type:TEXT;" json:"stock_name,omitempty"`
	ContactName  string `gorm:"column:contact_name;type:TEXT;" json:"contact_name,omitempty"`
	ContactPhone string `gorm:"column:contact_phone;type:TEXT;" json:"contact_phone,omitempty"`
	ContactEmail string `gorm:"column:contact_email;type:TEXT;" json:"contact_email,omitempty"`
	Owner        string `gorm:"column:owner;type:TEXT;" json:"owner,omitempty"`
	
	// 交貨資訊
	ExpectedDeliveryPeriod string     `gorm:"column:expected_delivery_period;type:TEXT;" json:"expected_delivery_period,omitempty"`
	ExpectedDeliveryDate   *time.Time `gorm:"column:expected_delivery_date;type:DATE;" json:"expected_delivery_date,omitempty"`
	ExpectedContractPeriod string     `gorm:"column:expected_contract_period;type:TEXT;" json:"expected_contract_period,omitempty"`
	ContractStartDate      *time.Time `gorm:"column:contract_start_date;type:DATE;" json:"contract_start_date,omitempty"`
	ContractEndDate        *time.Time `gorm:"column:contract_end_date;type:DATE;" json:"contract_end_date,omitempty"`
	DeliveryAddress        string     `gorm:"column:delivery_address;type:TEXT;" json:"delivery_address,omitempty"`
	SpecialRequirements    string     `gorm:"column:special_requirements;type:TEXT;" json:"special_requirements,omitempty"`
	
	// 備註
	Remark string `gorm:"column:remark;type:TEXT;" json:"remark,omitempty"`
	
	// 時間戳記
	CreatedTime time.Time  `gorm:"column:created_time;type:TIMESTAMP;" json:"created_time"`
	UpdatedTime *time.Time `gorm:"column:updated_time;type:TIMESTAMP;" json:"updated_time,omitempty"`
}

// Base struct is corresponding to table structure file
type Base struct {
	// 現貨編號
	StockID string `json:"stock_id,omitempty"`
	
	// 基本資訊
	StockName    string `json:"stock_name,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	Owner        string `json:"owner,omitempty"`
	
	// 交貨資訊
	ExpectedDeliveryPeriod string     `json:"expected_delivery_period,omitempty"`
	ExpectedDeliveryDate   *time.Time `json:"expected_delivery_date,omitempty"`
	ExpectedContractPeriod string     `json:"expected_contract_period,omitempty"`
	ContractStartDate      *time.Time `json:"contract_start_date,omitempty"`
	ContractEndDate        *time.Time `json:"contract_end_date,omitempty"`
	DeliveryAddress        string     `json:"delivery_address,omitempty"`
	SpecialRequirements    string     `json:"special_requirements,omitempty"`
	
	
	// 備註
	Remark string `json:"remark,omitempty"`
	
	// 時間戳記
	CreatedTime time.Time  `json:"created_time"`
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
}

// Single return structure file
type Single struct {
	// 現貨編號
	StockID string `json:"stock_id,omitempty"`
	
	// 基本資訊
	StockName    string `json:"stock_name,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	Owner        string `json:"owner,omitempty"`
	
	// 交貨資訊
	ExpectedDeliveryPeriod string     `json:"expected_delivery_period,omitempty"`
	ExpectedDeliveryDate   *time.Time `json:"expected_delivery_date,omitempty"`
	ExpectedContractPeriod string     `json:"expected_contract_period,omitempty"`
	ContractStartDate      *time.Time `json:"contract_start_date,omitempty"`
	ContractEndDate        *time.Time `json:"contract_end_date,omitempty"`
	DeliveryAddress        string     `json:"delivery_address,omitempty"`
	SpecialRequirements    string     `json:"special_requirements,omitempty"`
	
	
	// 備註
	Remark string `json:"remark,omitempty"`
	
	// 時間戳記
	CreatedTime time.Time  `json:"created_time"`
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
}

// Created struct is used to create
type Created struct {
	// 基本資訊
	StockName    string `json:"stock_name" binding:"required" validate:"required"`
	ContactName  string `json:"contact_name" binding:"required" validate:"required"`
	ContactPhone string `json:"contact_phone,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	Owner        string `json:"owner,omitempty"`
	
	// 交貨資訊
	ExpectedDeliveryPeriod string `json:"expected_delivery_period" binding:"required" validate:"required"`
	ExpectedDeliveryDate   string `json:"expected_delivery_date" binding:"required" validate:"required"`
	ExpectedContractPeriod string `json:"expected_contract_period" binding:"required" validate:"required"`
	ContractStartDate      string `json:"contract_start_date,omitempty"`
	ContractEndDate        string `json:"contract_end_date,omitempty"`
	DeliveryAddress        string `json:"delivery_address,omitempty"`
	SpecialRequirements    string `json:"special_requirements,omitempty"`
	
	
	// 備註
	Remark string `json:"remark,omitempty"`
}

// Field is structure file for search
type Field struct {
	// 現貨編號
	StockID      string  `json:"stock_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	StockName    *string `json:"stock_name,omitempty" form:"stock_name"`
	ContactName  *string `json:"contact_name,omitempty" form:"contact_name"`
}

// Fields is the searched structure file (including pagination)
type Fields struct {
	Field
	model.InPage
}

// List is multiple return structure files
type List struct {
	Stocks []*struct {
		// 現貨編號
		StockID string `json:"stock_id,omitempty"`
		
		// 基本資訊
		StockName    string `json:"stock_name,omitempty"`
		ContactName  string `json:"contact_name,omitempty"`
		ContactPhone string `json:"contact_phone,omitempty"`
		ContactEmail string `json:"contact_email,omitempty"`
		Owner        string `json:"owner,omitempty"`
		
		// 交貨資訊
		ExpectedDeliveryPeriod string     `json:"expected_delivery_period,omitempty"`
		ExpectedDeliveryDate   *time.Time `json:"expected_delivery_date,omitempty"`
		ExpectedContractPeriod string     `json:"expected_contract_period,omitempty"`
		ContractStartDate      *time.Time `json:"contract_start_date,omitempty"`
		ContractEndDate        *time.Time `json:"contract_end_date,omitempty"`
		DeliveryAddress        string     `json:"delivery_address,omitempty"`
		SpecialRequirements    string     `json:"special_requirements,omitempty"`
		
		
		// 備註
		Remark string `json:"remark,omitempty"`
		
		// 時間戳記
		CreatedTime time.Time  `json:"created_time"`
		UpdatedTime *time.Time `json:"updated_time,omitempty"`
	} `json:"stocks"`
	model.OutPage
}

// Updated struct is used to update
type Updated struct {
	// 現貨編號
	StockID string `json:"stock_id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	
	// 基本資訊
	StockName    string `json:"stock_name,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	Owner        string `json:"owner,omitempty"`
	
	// 交貨資訊
	ExpectedDeliveryPeriod string `json:"expected_delivery_period,omitempty"`
	ExpectedDeliveryDate   string `json:"expected_delivery_date,omitempty"`
	ExpectedContractPeriod string `json:"expected_contract_period,omitempty"`
	ContractStartDate      string `json:"contract_start_date,omitempty"`
	ContractEndDate        string `json:"contract_end_date,omitempty"`
	DeliveryAddress        string `json:"delivery_address,omitempty"`
	SpecialRequirements    string `json:"special_requirements,omitempty"`
	
	
	// 備註
	Remark string `json:"remark,omitempty"`
}

// TableName sets the insert table name for this struct type
func (a *Table) TableName() string {
	return "stocks"
}
