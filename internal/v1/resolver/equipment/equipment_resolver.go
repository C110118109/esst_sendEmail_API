package equipment

import (
	"encoding/json"
	"errors"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/linebot"
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

	// 設備建立完成後,發送第一階段 LINE 通知
	go r.sendProjectStep1LineNotification(input.ProjectID)

	return code.GetCodeMessage(code.Successful, "Batch created successfully")
}

// sendProjectStep1LineNotification 發送第一階段 LINE 通知
func (r *resolver) sendProjectStep1LineNotification(projectID string) {
	log.Info("Preparing to send step1 LINE notification for project:", projectID)

	// 查詢專案資訊
	projectField := &Field{ProjectID: projectID}
	project, err := r.ProjectService.GetByID(projectField.ToProjectField())
	if err != nil {
		log.Error("Failed to query project for LINE notification:", err)
		return
	}

	// 查詢設備清單
	equipments, err := r.EquipmentService.ListByProjectID(projectID)
	if err != nil {
		log.Error("Failed to query equipments for LINE notification:", err)
		return
	}

	// 轉換設備資料格式
	lineEquipments := make([]linebot.Equipment, 0)
	if equipments != nil {
		for _, eq := range equipments {
			lineEquipments = append(lineEquipments, linebot.Equipment{
				PartNumber:  eq.PartNumber,
				Quantity:    eq.Quantity,
				Description: eq.Description,
			})
		}
	}

	log.Info("Found", len(lineEquipments), "equipments for project", projectID)

	notificationData := &linebot.ProjectStep1Data{
		ProjectID:    project.ProjectID,
		ProjectName:  project.ProjectName,
		ContactName:  project.ContactName,
		ContactPhone: project.ContactPhone,
		ContactEmail: project.ContactEmail,
		Owner:        project.Owner,
		Remark:       project.Remark,
		Equipments:   lineEquipments,
		CreatedTime:  project.CreatedTime,
	}

	lineBotService := linebot.New()
	if err := lineBotService.SendProjectStep1Notification(notificationData); err != nil {
		log.Error("Failed to send step1 LINE notification:", err)
	} else {
		log.Info("Step1 LINE notification sent successfully for project:", projectID)
	}
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
