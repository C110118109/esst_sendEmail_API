package user

import (
	"encoding/json"
	"fmt"
	"time"

	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/users"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) Create(input *model.Created) (*model.Base, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 設定角色，預設為 user
	role := input.Role
	if role == "" {
		role = "user"
	}

	// 準備 table 結構
	table := &model.Table{
		ID:        util.GenerateUUID(),
		Username:  input.Username,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: time.Now(),
	}

	// 建立用戶
	err = s.Entity.Create(table)
	if err != nil {
		return nil, err
	}

	// 轉換為 Base 結構
	output := &model.Base{
		ID:        table.ID,
		Username:  table.Username,
		Role:      table.Role,
		CreatedAt: table.CreatedAt,
	}

	return output, nil
}

func (s *service) Authenticate(username, password string) (*model.Base, error) {
	// 查詢用戶
	field := &model.Field{Username: &username}
	user, err := s.Entity.GetByUsername(field)
	if err != nil {
		return nil, err
	}

	// 🔍 在比對前就印出來（非常重要）
	fmt.Println("🔍 Username:", username)
	fmt.Println("🔍 Input password:", password)
	fmt.Println("🔍 Stored hash:", user.Password)
	fmt.Println("🔍 Hash length:", len(user.Password))

	// 驗證密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("❌ Password compare failed:", err)
		return nil, err
	}

	// 轉換為 Base 結構
	output := &model.Base{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return output, nil
}

func (s *service) GetByID(input *model.Field) (*model.Base, error) {
	user, err := s.Entity.GetByID(input)
	if err != nil {
		return nil, err
	}

	// 轉換為 Base 結構
	output := &model.Base{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return output, nil
}

func (s *service) List(input *model.Fields) (int64, []*model.Base, error) {
	total, users, err := s.Entity.List(input)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	marshal, err := json.Marshal(users)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	var output []*model.Base
	err = json.Unmarshal(marshal, &output)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, output, nil
}

func (s *service) Update(input *model.Updated) error {
	id := input.ID
	user, err := s.Entity.GetByID(&model.Field{ID: &id})
	if err != nil {
		log.Error(err)
		return err
	}

	// 更新欄位
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	now := time.Now()
	user.UpdatedAt = &now

	return s.Entity.Update(user)
}

func (s *service) Delete(input *model.Field) error {
	return s.Entity.Delete(input)
}
