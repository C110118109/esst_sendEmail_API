package stock

import (
	"esst_sendEmail/internal/v1/resolver/stock"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Presenter interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type presenter struct {
	StockResolver stock.Resolver
}

func New(db *gorm.DB) Presenter {
	return &presenter{
		StockResolver: stock.New(db),
	}
}
