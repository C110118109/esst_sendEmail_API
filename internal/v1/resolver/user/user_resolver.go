package user

import (
	"encoding/json"
	"errors"

	"esst_sendEmail/internal/pkg/auth"
	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/users"
	"gorm.io/gorm"
)

func (r *resolver) Create(trx *gorm.DB, input *model.Created) interface{} {
	defer trx.Rollback()

	user, err := r.UserService.WithTrx(trx).Create(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()
	return code.GetCodeMessage(code.Successful, user.ID)
}

func (r *resolver) Login(input *model.Login) interface{} {
	// 驗證用戶
	user, err := r.UserService.Authenticate(input.Username, input.Password)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.JWTRejected, "Invalid username or password")
	}

	// 生成 JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.Error("Failed to generate token:", err)
		return code.GetCodeMessage(code.InternalServerError, "Failed to generate token")
	}

	// 將 token 加入回傳資料
	user.Token = token

	return code.GetCodeMessage(code.Successful, user)
}

func (r *resolver) GetByID(input *model.Field) interface{} {
	user, err := r.UserService.GetByID(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, user)
}

func (r *resolver) List(input *model.Fields) interface{} {
	output := &model.List{}
	output.Limit = input.Limit
	output.Page = input.Page

	quantity, users, err := r.UserService.List(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	usersByte, err := json.Marshal(users)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	output.Pages = util.Pagination(quantity, output.Limit)
	err = json.Unmarshal(usersByte, &output.Users)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, output)
}

func (r *resolver) Update(input *model.Updated) interface{} {
	// 驗證用戶是否存在
	id := input.ID
	_, err := r.UserService.GetByID(&model.Field{ID: &id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.UserService.Update(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, input.ID)
}

func (r *resolver) Delete(input *model.Field) interface{} {
	// 驗證用戶是否存在
	_, err := r.UserService.GetByID(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.UserService.Delete(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, "Delete ok!")
}
