package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"esst_sendEmail/internal/v1/router/equipment"
	"esst_sendEmail/internal/v1/router/project"
	"esst_sendEmail/internal/v1/router/stock"
	"esst_sendEmail/internal/v1/router/stock_equipment"

	"esst_sendEmail/internal/v1/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS for cross-origin requests
	router.Use(middleware.CORSMiddleware())

	// 註冊路由
	router = project.GetRoute(router, db)
	router = equipment.GetRoute(router, db)
	router = stock.GetRoute(router, db)
	router = stock_equipment.GetRoute(router, db)

	log.Fatal(http.ListenAndServe(":8080", router))
}
