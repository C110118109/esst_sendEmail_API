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
	projectID := ctx.Param("projectID")
	input := &projects.Field{}
	input.ProjectID = projectID
	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.FormatError, err.Error()))

		return
	}

	codeMessage := p.ProjectResolver.GetByID(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Update(ctx *gin.Context) {
	// Todo 將UUID改成登入的使用者
	//updatedBy := util.GenerateUUID()
	projectID := ctx.Param("projectID")
	input := &projects.Updated{}
	input.ProjectID = projectID
	//input.UpdatedBy = util.PointerString(updatedBy)
	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.FormatError, err.Error()))

		return
	}

	codeMessage := p.ProjectResolver.Update(input)
	ctx.JSON(http.StatusOK, codeMessage)
}

func (p *presenter) Delete(ctx *gin.Context) {
	// Todo 將UUID改成登入的使用者
	//updatedBy := util.GenerateUUID()
	projectID := ctx.Param("projectID")
	input := &projects.Updated{}
	input.ProjectID = projectID
	//input.UpdatedBy = util.PointerString(updatedBy)
	if err := ctx.ShouldBindJSON(input); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusOK, code.GetCodeMessage(code.FormatError, err.Error()))

		return
	}

	codeMessage := p.ProjectResolver.Delete(input)
	ctx.JSON(http.StatusOK, codeMessage)
}
