package user

import (
	"esst_sendEmail/internal/v1/resolver/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Presenter interface {
	Create(ctx *gin.Context)
	RequestVerificationCode(ctx *gin.Context) // 新增:請求驗證碼
	VerifyAndLogin(ctx *gin.Context)          // 新增:驗證碼登入
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	List(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type presenter struct {
	UserResolver user.Resolver
}

func New(db *gorm.DB) Presenter {
	return &presenter{
		UserResolver: user.New(db),
	}
}
