package users

import (
	"time"

	model "esst_sendEmail/internal/v1/structure"
)

// Table 資料表結構
type Table struct {
	ID        string     `gorm:"primaryKey;uuid_generate_v4();column:id;type:uuid;" json:"id,omitempty"`
	Username  string     `gorm:"column:username;type:TEXT;unique;" json:"username,omitempty"`
	Email     string     `gorm:"column:email;type:TEXT;" json:"email,omitempty"` // 新增 email 欄位
	Password  string     `gorm:"column:password;type:TEXT;" json:"-"`
	Role      string     `gorm:"column:role;type:TEXT;default:'user';" json:"role,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;type:TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:TIMESTAMP;" json:"updated_at,omitempty"`
}

// Base 基礎結構
type Base struct {
	ID        string     `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Email     string     `json:"email,omitempty"` // 新增 email 欄位
	Role      string     `json:"role,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Token     string     `json:"token,omitempty"` // 登入時返回 token
}

// Created 建立用戶
type Created struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"` // 新增 email 欄位
	Password string `json:"password" binding:"required" validate:"required"`
	Role     string `json:"role,omitempty"` // 預設為 user
}

// Login 登入
type Login struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// Field 查詢條件
type Field struct {
	ID       *string `json:"id,omitempty" binding:"omitempty,uuid4"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"` // 新增 email 查詢欄位
}

// Fields 多筆查詢
type Fields struct {
	Field
	model.InPage
}

// List 列表結構
type List struct {
	Users []*struct {
		ID        string     `json:"id,omitempty"`
		Username  string     `json:"username,omitempty"`
		Email     string     `json:"email,omitempty"` // 新增 email 欄位
		Role      string     `json:"role,omitempty"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	} `json:"users"`
	model.OutPage
}

// Updated 更新用戶
type Updated struct {
	ID       string `json:"id,omitempty" binding:"omitempty,uuid4" swaggerignore:"true"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"` // 新增 email 欄位
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

// TableName 設定資料表名稱
func (t *Table) TableName() string {
	return "users"
}
