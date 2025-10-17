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
	//input.IsDeleted = util.PointerBool(false)
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
	// department, err := r.DepartmentService.GetByID(&departmentModel.Field{DepartmentID: input.DepartmentID,
	// 	IsDeleted: util.PointerBool(false)})
	project, err := r.ProjectService.GetByID(&model.Field{ProjectID: input.ProjectID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.ProjectService.Update(input)
	if err != nil {
		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, project.ProjectID)
}

func (r *resolver) Delete(input *model.Updated) interface{} {
	// _, err := r.DepartmentService.GetByID(&departmentModel.Field{DepartmentID: input.DepartmentID,
	// 	IsDeleted: util.PointerBool(false)})
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
