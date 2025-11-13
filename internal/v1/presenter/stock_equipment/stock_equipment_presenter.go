package stock_equipment

import (
	"net/http"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	preset "esst_sendEmail/internal/v1/presenter"
	"esst_sendEmail/internal/v1/structure/stock_equipments"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (p *presenter) Create(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &stock_equipments.Created{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.StockEquipmentResolver.Create(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) CreateBatch(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &stock_equipments.BatchCreated{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.StockEquipmentResolver.CreateBatch(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) List(ctx *gin.Context) {
	input := &stock_equipments.Fields{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	if input.Limit == 0 || input.Limit > preset.DefaultLimit {
		input.Limit = preset.DefaultLimit
	}

	codeMessage := p.StockEquipmentResolver.List(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) ListByStockID(ctx *gin.Context) {
	stockID := ctx.Param("stockId")

	if stockID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Stock ID is required"))
		return
	}

	codeMessage := p.StockEquipmentResolver.ListByStockID(stockID)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) GetByID(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &stock_equipments.Field{}
	input.StockEquipmentID = equipmentID

	if equipmentID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Equipment ID is required"))
		return
	}

	codeMessage := p.StockEquipmentResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Update(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &stock_equipments.Updated{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	input.StockEquipmentID = equipmentID

	codeMessage := p.StockEquipmentResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Delete(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &stock_equipments.Updated{}
	input.StockEquipmentID = equipmentID

	if equipmentID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Equipment ID is required"))
		return
	}

	codeMessage := p.StockEquipmentResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
