package equipment

import (
	"esst_sendEmail/internal/v1/resolver/equipment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Presenter interface {
	Create(ctx *gin.Context)
	CreateBatch(ctx *gin.Context)
	List(ctx *gin.Context)
	ListByProjectID(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type presenter struct {
	EquipmentResolver equipment.Resolver
}

func New(db *gorm.DB) Presenter {
	return &presenter{
		EquipmentResolver: equipment.New(db),
	}
}
