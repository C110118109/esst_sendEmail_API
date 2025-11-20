package stock_equipment

import (
	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/presenter/stock_equipment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := stock_equipment.New(db)
	v10 := route.Group("authority").Group("v1.0").Group("stock-equipments")
	v10.Use(middleware.JWTMiddleware()) // 加上 JWT 驗證
	{
		// 單筆建立現貨設備
		v10.POST("", middleware.Transaction(db), controller.Create)
		// 批次建立現貨設備
		v10.POST("/batch", middleware.Transaction(db), controller.CreateBatch)
		// 獲取現貨設備列表
		v10.GET("", controller.List)
		// 根據現貨ID獲取設備列表
		v10.GET("/stock/:stockId", controller.ListByStockID)
		// 根據設備ID獲取單筆設備
		v10.GET("/:equipmentId", controller.GetByID)
		// 更新設備
		v10.PATCH("/:equipmentId", controller.Update)
		// 刪除設備
		v10.DELETE("/:equipmentId", controller.Delete)
	}
	return route
}
