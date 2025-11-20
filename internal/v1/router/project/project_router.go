package project

import (
	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/presenter/project"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := project.New(db)
	v10 := route.Group("authority").Group("v1.0").Group("projects")
	v10.Use(middleware.JWTMiddleware()) // 加上 JWT 驗證
	{
		// 建立專案（第一階段）
		v10.POST("", middleware.Transaction(db), controller.Create)
		// 查詢專案列表
		v10.GET("", controller.List)
		// 查詢單一專案
		v10.GET("/:projectId", controller.GetByID)
		// 刪除專案
		v10.DELETE("/:projectId", controller.Delete)
		// 更新專案（包含第二階段）
		v10.PATCH("/:projectId", controller.Update)
	}

	return route
}
