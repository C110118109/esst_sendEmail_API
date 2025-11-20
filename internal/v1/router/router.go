package router

import (
	"github.com/gin-gonic/gin"
)

func Default() *gin.Engine {

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}
