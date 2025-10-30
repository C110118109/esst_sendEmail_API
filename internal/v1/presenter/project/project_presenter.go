package project

import (
	"net/http"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	preset "esst_sendEmail/internal/v1/presenter"
	"esst_sendEmail/internal/v1/structure/projects"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (p *presenter) Create(ctx *gin.Context) {
	trx := ctx.MustGet("db_trx").(*gorm.DB)
	input := &projects.Created{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.ProjectResolver.Create(trx, input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) List(ctx *gin.Context) {
	input := &projects.Fields{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	if input.Limit == 0 || input.Limit > preset.DefaultLimit {
		input.Limit = preset.DefaultLimit
	}

	codeMessage := p.ProjectResolver.List(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) GetByID(ctx *gin.Context) {
	// 修正：URL 參數名稱改為 projectId（與 router 定義一致）
	projectId := ctx.Param("projectId")
	input := &projects.Field{}
	input.ProjectID = projectId

	// GetByID 不需要從 body 取得資料，直接使用 URL 參數
	codeMessage := p.ProjectResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Update(ctx *gin.Context) {
	// 修正：URL 參數名稱改為 projectId（與 router 定義一致）
	projectId := ctx.Param("projectId")
	input := &projects.Updated{}
	input.ProjectID = projectId

	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, code.GetCodeMessage(code.FormatError, err.Error()))
		return
	}

	codeMessage := p.ProjectResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Delete(ctx *gin.Context) {
	// 修正：URL 參數名稱改為 projectId（與 router 定義一致）
	projectId := ctx.Param("projectId")
	input := &projects.Updated{}
	input.ProjectID = projectId

	// Delete 操作不需要從 body 取得額外資料
	codeMessage := p.ProjectResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
