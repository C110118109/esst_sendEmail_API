package user

import (
	"net/http"
	"time"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	preset "esst_sendEmail/internal/v1/presenter"
	"esst_sendEmail/internal/v1/structure/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create 建立用戶 (僅限管理員)
func (p *presenter) Create(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &users.Created{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.UserResolver.Create(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

// Login 登入
func (p *presenter) Login(ctx *gin.Context) {
	input := &users.Login{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	// 執行登入
	result := p.UserResolver.Login(input)

	// 檢查登入是否成功
	successMsg, ok := result.(*code.SuccessfulMessage)
	if ok && successMsg.Code == code.Successful {
		// 從成功訊息中獲取用戶資料
		userData, ok := successMsg.Body.(*users.Base)
		if !ok {
			log.Error("Failed to convert user data")
			ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "Internal server error"))
			return
		}

		// 設定 token cookie
		ctx.SetCookie(
			"token",
			userData.Token,
			int(24*time.Hour.Seconds()), // 24 小時
			"/",
			"",
			false, // 開發環境設為 false，生產環境改為 true
			true,
		)

		// 返回帶有 token 的回應
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.Successful, userData))
	} else {
		// 登入失敗，直接返回錯誤結果
		ctx.JSON(http.StatusOK, result)
	}
}

// Logout 登出
func (p *presenter) Logout(ctx *gin.Context) {
	// 清除 token cookie
	ctx.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, code.GetCodeMessage(code.Successful, "Logged out successfully"))
}

// GetByID 取得單一用戶資訊
func (p *presenter) GetByID(ctx *gin.Context) {
	userId := ctx.Param("userId")
	input := &users.Field{}
	input.ID = &userId

	codeMessage := p.UserResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

// List 取得用戶列表 (僅限管理員)
func (p *presenter) List(ctx *gin.Context) {
	input := &users.Fields{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	if input.Limit == 0 || input.Limit > preset.DefaultLimit {
		input.Limit = preset.DefaultLimit
	}

	codeMessage := p.UserResolver.List(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

// Update 更新用戶 (僅限管理員)
func (p *presenter) Update(ctx *gin.Context) {
	userId := ctx.Param("userId")
	input := &users.Updated{}
	input.ID = userId

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.UserResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

// Delete 刪除用戶 (僅限管理員)
func (p *presenter) Delete(ctx *gin.Context) {
	userId := ctx.Param("userId")
	input := &users.Field{}
	input.ID = &userId

	codeMessage := p.UserResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
