package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/router/equipment"
	"esst_sendEmail/internal/v1/router/project"
	"esst_sendEmail/internal/v1/router/stock"
	"esst_sendEmail/internal/v1/router/stock_equipment"
	"esst_sendEmail/internal/v1/router/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 資料庫連線設定
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

	// 初始化 Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 設定 CORS
	router.Use(middleware.CORSMiddleware())

	// 註冊路由
	// 1. 用戶認證路由（包含公開的 login 和需要 JWT 的管理功能）
	router = user.GetRoute(router, db)

	// 2. 專案管理路由（需要 JWT 驗證）
	router = project.GetRoute(router, db)

	// 3. 設備管理路由（需要 JWT 驗證）
	router = equipment.GetRoute(router, db)

	// 4. 現貨管理路由（需要 JWT 驗證）
	router = stock.GetRoute(router, db)

	// 5. 現貨設備管理路由（需要 JWT 驗證）
	router = stock_equipment.GetRoute(router, db)

	// 啟動服務器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Printf("預設管理員帳號: admin / 密碼: admin123")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
