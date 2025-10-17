package equipment

import (
	"net/http"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	preset "esst_sendEmail/internal/v1/presenter"
	"esst_sendEmail/internal/v1/structure/equipments"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (p *presenter) Create(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &equipments.Created{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.EquipmentResolver.Create(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) CreateBatch(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &equipments.BatchCreated{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.EquipmentResolver.CreateBatch(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) List(ctx *gin.Context) {
	input := &equipments.Fields{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	if input.Limit == 0 || input.Limit > preset.DefaultLimit {
		input.Limit = preset.DefaultLimit
	}

	codeMessage := p.EquipmentResolver.List(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) ListByProjectID(ctx *gin.Context) {
	projectID := ctx.Param("projectId")

	if projectID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Project ID is required"))
		return
	}

	codeMessage := p.EquipmentResolver.ListByProjectID(projectID)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) GetByID(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &equipments.Field{}
	input.EquipmentID = equipmentID

	// 修正: GET 請求不應該使用 ShouldBindJSON,直接使用從 URL 取得的參數
	if equipmentID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Equipment ID is required"))
		return
	}

	codeMessage := p.EquipmentResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Update(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &equipments.Updated{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	input.EquipmentID = equipmentID

	codeMessage := p.EquipmentResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Delete(ctx *gin.Context) {
	equipmentID := ctx.Param("equipmentId")
	input := &equipments.Updated{}
	input.EquipmentID = equipmentID

	// DELETE 請求通常不需要 body,直接使用 URL 參數
	if equipmentID == "" {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, "Equipment ID is required"))
		return
	}

	codeMessage := p.EquipmentResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
