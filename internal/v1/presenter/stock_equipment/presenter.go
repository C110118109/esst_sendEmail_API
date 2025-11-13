package stock_equipment

import (
	"esst_sendEmail/internal/v1/resolver/stock_equipment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Presenter interface {
	Create(ctx *gin.Context)
	CreateBatch(ctx *gin.Context)
	List(ctx *gin.Context)
	ListByStockID(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type presenter struct {
	StockEquipmentResolver stock_equipment.Resolver
}

func New(db *gorm.DB) Presenter {
	return &presenter{
		StockEquipmentResolver: stock_equipment.New(db),
	}
}
