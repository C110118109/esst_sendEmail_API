package equipment

import (
	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/presenter/equipment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := equipment.New(db)
	v10 := route.Group("authority").Group("v1.0").Group("equipments")
	{
		// 單筆建立設備
		v10.POST("", middleware.Transaction(db), controller.Create)
		// 批次建立設備
		v10.POST("/batch", middleware.Transaction(db), controller.CreateBatch)
		// 獲取設備列表
		v10.GET("", controller.List)
		// 根據專案ID獲取設備列表
		v10.GET("/project/:projectId", controller.ListByProjectID)
		// 根據設備ID獲取單筆設備
		v10.GET("/:equipmentId", controller.GetByID)
		// 更新設備
		v10.PATCH("/:equipmentId", controller.Update)
		// 刪除設備
		v10.DELETE("/:equipmentId", controller.Delete)
	}
	return route
}
