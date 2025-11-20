package stock

import (
	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/presenter/stock"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := stock.New(db)
	v10 := route.Group("authority").Group("v1.0").Group("stocks")
	v10.Use(middleware.JWTMiddleware()) // 加上 JWT 驗證
	{
		// 建立現貨報備
		v10.POST("", middleware.Transaction(db), controller.Create)
		// 查詢現貨報備列表
		v10.GET("", controller.List)
		// 查詢單一現貨報備
		v10.GET("/:stockId", controller.GetByID)
		// 刪除現貨報備
		v10.DELETE("/:stockId", controller.Delete)
		// 更新現貨報備
		v10.PATCH("/:stockId", controller.Update)
	}

	return route
}
