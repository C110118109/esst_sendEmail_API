package project

import (
	"encoding/json"
	"errors"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/linebot"
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

	// 注意: 第一階段 LINE 通知將在設備建立完成後發送
	// 見 equipment/equipment_resolver.go 的 CreateBatch 函數

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

	// 檢查是否為第二階段更新
	isStep2Update := false
	if input.ExpectedDeliveryPeriod != "" ||
		input.ExpectedDeliveryDate != "" ||
		input.ExpectedContractPeriod != "" ||
		input.ContractStartDate != "" ||
		input.ContractEndDate != "" ||
		input.DeliveryAddress != "" ||
		input.SpecialRequirements != "" {
		isStep2Update = true
		// 如果狀態還是 step1,則更新為 step2
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

	// 如果是第二階段更新,發送 LINE 通知
	if isStep2Update {
		go func() {
			// 查詢設備清單
			equipments, err := r.EquipmentService.ListByProjectID(input.ProjectID)
			if err != nil {
				log.Error("Failed to query equipments for step2 LINE notification:", err)
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

			log.Info("Found", len(lineEquipments), "equipments for step2 LINE notification, project:", input.ProjectID)

			// 格式化日期
			formatDate := func(dateStr string) string {
				if dateStr == "" {
					return "-"
				}
				return dateStr
			}

			notificationData := &linebot.ProjectStep2Data{
				ProjectID:              input.ProjectID,
				ProjectName:            project.ProjectName,
				ContactName:            project.ContactName,
				ExpectedDeliveryPeriod: input.ExpectedDeliveryPeriod,
				ExpectedDeliveryDate:   formatDate(input.ExpectedDeliveryDate),
				ExpectedContractPeriod: input.ExpectedContractPeriod,
				ContractStartDate:      formatDate(input.ContractStartDate),
				ContractEndDate:        formatDate(input.ContractEndDate),
				DeliveryAddress:        input.DeliveryAddress,
				SpecialRequirements:    input.SpecialRequirements,
				Equipments:             lineEquipments,
				UpdatedTime:            util.NowToUTC(),
			}

			lineBotService := linebot.New()
			if err := lineBotService.SendProjectStep2Notification(notificationData); err != nil {
				log.Error("Failed to send step2 LINE notification:", err)
			} else {
				log.Info("Step2 LINE notification sent successfully for project:", input.ProjectID)
			}
		}()
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
