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
	{
		v10.POST("", middleware.Transaction(db), controller.Create)
		v10.GET("", controller.List)
		v10.GET("/:projectId", controller.GetByID)
		v10.DELETE("/:projectId", controller.Delete)
		v10.PATCH("/:projectId", controller.Update)
	}

	return route
}
