package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost",
			"http://localhost:80",
			"http://localhost:5500",
			"http://localhost:8080",
			"http://127.0.0.1",
			"http://127.0.0.1:80",
			"http://127.0.0.1:5500",
			"http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true, // Critical for cookies/sessions
		MaxAge:           12 * time.Hour,
	})
}
