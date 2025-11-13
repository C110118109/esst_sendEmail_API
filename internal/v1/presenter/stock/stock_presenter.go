package stock

import (
	"net/http"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	preset "esst_sendEmail/internal/v1/presenter"
	"esst_sendEmail/internal/v1/structure/stocks"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (p *presenter) Create(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &stocks.Created{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.StockResolver.Create(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) List(ctx *gin.Context) {
	input := &stocks.Fields{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	if input.Limit == 0 || input.Limit > preset.DefaultLimit {
		input.Limit = preset.DefaultLimit
	}

	codeMessage := p.StockResolver.List(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) GetByID(ctx *gin.Context) {
	stockId := ctx.Param("stockId")
	input := &stocks.Field{}
	input.StockID = stockId

	codeMessage := p.StockResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Update(ctx *gin.Context) {
	stockId := ctx.Param("stockId")
	input := &stocks.Updated{}
	input.StockID = stockId

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.StockResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Delete(ctx *gin.Context) {
	stockId := ctx.Param("stockId")
	input := &stocks.Updated{}
	input.StockID = stockId

	codeMessage := p.StockResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
