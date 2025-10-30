package project

import (
	"encoding/json"
	"errors"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/projects"

	"gorm.io/gorm"
)

func (r *resolver) Create(trx *gorm.DB, input *model.Created) interface{} {
	defer trx.Rollback()

	project, err := r.ProjectService.WithTrx(trx).Create(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()
	return code.GetCodeMessage(code.Successful, project.ProjectID)
}

func (r *resolver) List(input *model.Fields) interface{} {
	output := &model.List{}
	output.Limit = input.Limit
	output.Page = input.Page

	quantity, projects, err := r.ProjectService.List(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	projectsByte, err := json.Marshal(projects)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	output.Pages = util.Pagination(quantity, output.Limit)
	err = json.Unmarshal(projectsByte, &output.Projects)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, output)
}

func (r *resolver) GetByID(input *model.Field) interface{} {
	base, err := r.ProjectService.GetByID(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	frontProject := &model.Single{}
	projectsByte, _ := json.Marshal(base)
	err = json.Unmarshal(projectsByte, &frontProject)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, frontProject)
}

func (r *resolver) Update(input *model.Updated) interface{} {
	// 驗證專案是否存在
	project, err := r.ProjectService.GetByID(&model.Field{ProjectID: input.ProjectID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	// 如果有第二階段欄位更新，自動將狀態改為 step2
	if input.ExpectedDeliveryPeriod != "" ||
		input.ExpectedDeliveryDate != "" ||
		input.ExpectedContractPeriod != "" ||
		input.ContractStartDate != "" ||
		input.ContractEndDate != "" ||
		input.DeliveryAddress != "" ||
		input.SpecialRequirements != "" {
		// 如果狀態還是 step1，則更新為 step2
		if project.Status == "step1" && input.Status == "" {
			input.Status = "step2"
		}
	}

	// 執行更新
	err = r.ProjectService.Update(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, project.ProjectID)
}

func (r *resolver) Delete(input *model.Updated) interface{} {
	// 驗證專案是否存在
	_, err := r.ProjectService.GetByID(&model.Field{ProjectID: input.ProjectID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.ProjectService.Delete(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, "Delete ok!")
}
