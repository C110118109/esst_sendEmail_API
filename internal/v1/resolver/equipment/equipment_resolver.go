package equipment

import (
	"encoding/json"
	"errors"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/equipments"

	"gorm.io/gorm"
)

func (r *resolver) Create(trx *gorm.DB, input *model.Created) interface{} {
	defer trx.Rollback()

	equipment, err := r.EquipmentService.WithTrx(trx).Create(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()
	return code.GetCodeMessage(code.Successful, equipment.EquipmentID)
}

func (r *resolver) CreateBatch(trx *gorm.DB, input *model.BatchCreated) interface{} {
	defer trx.Rollback()

	err := r.EquipmentService.WithTrx(trx).CreateBatch(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()
	return code.GetCodeMessage(code.Successful, "Batch created successfully")
}

func (r *resolver) List(input *model.Fields) interface{} {
	output := &model.List{}
	output.Limit = input.Limit
	output.Page = input.Page

	quantity, equipments, err := r.EquipmentService.List(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	equipmentsByte, err := json.Marshal(equipments)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	output.Pages = util.Pagination(quantity, output.Limit)
	err = json.Unmarshal(equipmentsByte, &output.Equipments)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, output)
}

func (r *resolver) ListByProjectID(projectID string) interface{} {
	equipments, err := r.EquipmentService.ListByProjectID(projectID)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, equipments)
}

func (r *resolver) GetByID(input *model.Field) interface{} {
	base, err := r.EquipmentService.GetByID(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	frontEquipment := &model.Single{}
	equipmentsByte, _ := json.Marshal(base)
	err = json.Unmarshal(equipmentsByte, &frontEquipment)
	if err != nil {
		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, frontEquipment)
}

func (r *resolver) Update(input *model.Updated) interface{} {
	equipment, err := r.EquipmentService.GetByID(&model.Field{EquipmentID: input.EquipmentID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.EquipmentService.Update(input)
	if err != nil {
		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, equipment.EquipmentID)
}

func (r *resolver) Delete(input *model.Updated) interface{} {
	_, err := r.EquipmentService.GetByID(&model.Field{EquipmentID: input.EquipmentID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.EquipmentService.Delete(input)
	if err != nil {
		log.Error(err)

		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, "Delete ok!")
}
