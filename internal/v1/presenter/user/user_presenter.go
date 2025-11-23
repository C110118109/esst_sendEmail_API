package user

import (
	"net/http"
	"time"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/mail"
	"esst_sendEmail/internal/pkg/verification"
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

// RequestVerificationCode 請求驗證碼 - 第一步驗證帳號密碼
func (p *presenter) RequestVerificationCode(ctx *gin.Context) {
	input := &users.Login{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	// 先驗證帳號密碼
	result := p.UserResolver.ValidateCredentials(input.Username, input.Password)

	// 檢查驗證是否成功
	successMsg, ok := result.(*code.SuccessfulMessage)
	if !ok || successMsg.Code != code.Successful {
		ctx.JSON(http.StatusOK, result)
		return
	}

	// 獲取使用者資料
	userData, ok := successMsg.Body.(*users.Base)
	if !ok {
		log.Error("Failed to convert user data")
		ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "Internal server error"))
		return
	}

	// 檢查使用者是否有設定 email
	if userData.Email == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "使用者未設定電子信箱，無法發送驗證碼"))
		return
	}

	// 生成驗證碼
	verificationService := verification.New()
	verificationCode, err := verificationService.GenerateCode(userData.Email)
	if err != nil {
		log.Error("Failed to generate verification code:", err)
		ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "Failed to generate verification code"))
		return
	}

	// 發送驗證碼郵件
	emailService := mail.New()
	err = emailService.SendVerificationCode(userData.Email, verificationCode, userData.Username)
	if err != nil {
		log.Error("Failed to send verification code email:", err)
		ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "無法發送驗證碼郵件"))
		return
	}

	log.Info("Verification code sent to:", userData.Email)

	// 返回成功訊息(不返回驗證碼)
	ctx.JSON(http.StatusOK, code.GetCodeMessage(code.Successful, gin.H{
		"message": "驗證碼已發送至您的信箱",
		"email":   maskEmail(userData.Email), // 遮罩部分信箱
	}))
}

// VerifyAndLogin 驗證驗證碼並登入 - 第二步驗證驗證碼
func (p *presenter) VerifyAndLogin(ctx *gin.Context) {
	type VerifyInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required,len=6"`
	}

	input := &VerifyInput{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	// 再次驗證帳號密碼(避免重放攻擊)
	result := p.UserResolver.ValidateCredentials(input.Username, input.Password)

	successMsg, ok := result.(*code.SuccessfulMessage)
	if !ok || successMsg.Code != code.Successful {
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.JWTRejected, "帳號或密碼錯誤"))
		return
	}

	userData, ok := successMsg.Body.(*users.Base)
	if !ok {
		log.Error("Failed to convert user data")
		ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "Internal server error"))
		return
	}

	// 驗證驗證碼
	verificationService := verification.New()
	if !verificationService.VerifyCode(userData.Email, input.Code) {
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.JWTRejected, "驗證碼錯誤或已過期"))
		return
	}

	// 驗證成功，執行登入
	loginResult := p.UserResolver.Login(&users.Login{
		Username: input.Username,
		Password: input.Password,
	})

	// 檢查登入是否成功
	loginSuccessMsg, ok := loginResult.(*code.SuccessfulMessage)
	if ok && loginSuccessMsg.Code == code.Successful {
		// 從成功訊息中獲取用戶資料
		loginUserData, ok := loginSuccessMsg.Body.(*users.Base)
		if !ok {
			log.Error("Failed to convert user data")
			ctx.JSON(http.StatusInternalServerError, code.GetCodeMessage(code.InternalServerError, "Internal server error"))
			return
		}

		// 設定 token cookie
		ctx.SetCookie(
			"token",
			loginUserData.Token,
			int(24*time.Hour.Seconds()), // 24 小時
			"/",
			"",
			false, // 開發環境設為 false，生產環境改為 true
			true,
		)

		// 返回帶有 token 的回應
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.Successful, loginUserData))
	} else {
		// 登入失敗，直接返回錯誤結果
		ctx.JSON(http.StatusOK, loginResult)
	}
}

// Login 原有的直接登入(保留作為備用)
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

// maskEmail 遮罩信箱地址
func maskEmail(email string) string {
	if email == "" {
		return ""
	}

	// 找到 @ 的位置
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}

	if atIndex < 0 {
		return email
	}

	// 保留前 2 個字元和 @ 後面的部分
	if atIndex <= 2 {
		return email
	}

	masked := email[:2]
	for i := 2; i < atIndex; i++ {
		masked += "*"
	}
	masked += email[atIndex:]

	return masked
}
