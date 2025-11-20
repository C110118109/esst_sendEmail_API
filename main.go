package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"esst_sendEmail/internal/v1/middleware"
	"esst_sendEmail/internal/v1/router/equipment"
	"esst_sendEmail/internal/v1/router/project"
	"esst_sendEmail/internal/v1/router/stock"
	"esst_sendEmail/internal/v1/router/stock_equipment"
	"esst_sendEmail/internal/v1/router/user"
	userModel "esst_sendEmail/internal/v1/structure/users"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 驗證必要的環境變數
	requiredEnvs := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", env)
		}
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

	// 檢查並建立預設管理員帳號
	initDefaultAdmin(db)

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

	log.Printf("✅ Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// initDefaultAdmin 檢查並建立預設管理員帳號
func initDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&userModel.Table{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		// 生成密碼 hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash default admin password:", err)
		}

		// 建立預設管理員
		adminUser := &userModel.Table{
			ID:        uuid.New().String(),
			Username:  "admin",
			Password:  string(hashedPassword),
			Role:      "admin",
			CreatedAt: time.Now(),
		}

		if err := db.Create(adminUser).Error; err != nil {
			log.Fatal("Failed to create default admin user:", err)
		}

		log.Println("⚠️  預設管理員已建立")
		log.Println("⚠️  帳號: admin")
		log.Println("⚠️  密碼: admin123")
		log.Println("⚠️  請立即登入並修改密碼!")
	}
}
