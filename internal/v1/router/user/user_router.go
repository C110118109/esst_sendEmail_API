package user

import (
	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/presenter/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := user.New(db)

	// 公開路由 - 登入相關
	route.POST("/auth/request-code", controller.RequestVerificationCode) // 新增:請求驗證碼
	route.POST("/auth/verify-login", controller.VerifyAndLogin)          // 新增:驗證碼登入
	route.POST("/auth/login", controller.Login)                          // 保留原有的直接登入
	route.POST("/auth/logout", controller.Logout)

	// 需要身份驗證的路由
	auth := route.Group("/auth")
	auth.Use(middleware.JWTMiddleware())
	{
		// 取得當前用戶資訊
		auth.GET("/me", func(ctx *gin.Context) {
			userID, _ := ctx.Get("userID")
			username, _ := ctx.Get("username")
			role, _ := ctx.Get("role")

			ctx.JSON(200, gin.H{
				"id":       userID,
				"username": username,
				"role":     role,
			})
		})
	}

	// 管理員路由 - 用戶管理
	v10 := route.Group("authority").Group("v1.0").Group("users")
	v10.Use(middleware.JWTMiddleware(), middleware.AdminMiddleware())
	{
		// 建立用戶
		v10.POST("", middleware.Transaction(db), controller.Create)
		// 查詢用戶列表
		v10.GET("", controller.List)
		// 查詢單一用戶
		v10.GET("/:userId", controller.GetByID)
		// 更新用戶
		v10.PATCH("/:userId", controller.Update)
		// 刪除用戶
		v10.DELETE("/:userId", controller.Delete)
	}

	return route
}
