package stock_equipment

import (
	"encoding/json"
	"errors"
	"time"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/linebot"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/stock_equipments"

	"gorm.io/gorm"
)

func (r *resolver) Create(trx *gorm.DB, input *model.Created) interface{} {
	defer trx.Rollback()

	equipment, err := r.StockEquipmentService.WithTrx(trx).Create(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()
	return code.GetCodeMessage(code.Successful, equipment.StockEquipmentID)
}

func (r *resolver) CreateBatch(trx *gorm.DB, input *model.BatchCreated) interface{} {
	defer trx.Rollback()

	err := r.StockEquipmentService.WithTrx(trx).CreateBatch(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()

	// 設備建立完成後,發送現貨報備 LINE 通知
	go r.sendStockLineNotification(input.StockID)

	return code.GetCodeMessage(code.Successful, "Batch created successfully")
}

// sendStockLineNotification 發送現貨報備 LINE 通知
func (r *resolver) sendStockLineNotification(stockID string) {
	log.Info("Preparing to send stock LINE notification for stock:", stockID)

	// 查詢現貨資訊
	stockField := &Field{StockID: stockID}
	stock, err := r.StockService.GetByID(stockField.ToStockField())
	if err != nil {
		log.Error("Failed to query stock for LINE notification:", err)
		return
	}

	// 查詢設備清單
	equipments, err := r.StockEquipmentService.ListByStockID(stockID)
	if err != nil {
		log.Error("Failed to query stock equipments for LINE notification:", err)
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

	log.Info("Found", len(lineEquipments), "equipments for stock", stockID)

	// 格式化日期
	formatDate := func(datePtr *time.Time) string {
		if datePtr == nil {
			return "-"
		}
		return datePtr.Format("2006-01-02")
	}

	notificationData := &linebot.StockData{
		StockID:                stock.StockID,
		StockName:              stock.StockName,
		ContactName:            stock.ContactName,
		ContactPhone:           stock.ContactPhone,
		ContactEmail:           stock.ContactEmail,
		Owner:                  stock.Owner,
		ExpectedDeliveryPeriod: stock.ExpectedDeliveryPeriod,
		ExpectedDeliveryDate:   formatDate(stock.ExpectedDeliveryDate),
		ExpectedContractPeriod: stock.ExpectedContractPeriod,
		ContractStartDate:      formatDate(stock.ContractStartDate),
		ContractEndDate:        formatDate(stock.ContractEndDate),
		DeliveryAddress:        stock.DeliveryAddress,
		SpecialRequirements:    stock.SpecialRequirements,
		Remark:                 stock.Remark,
		Equipments:             lineEquipments,
		CreatedTime:            stock.CreatedTime,
	}

	lineBotService := linebot.New()
	if err := lineBotService.SendStockNotification(notificationData); err != nil {
		log.Error("Failed to send stock LINE notification:", err)
	} else {
		log.Info("Stock LINE notification sent successfully for stock:", stockID)
	}
}

func (r *resolver) List(input *model.Fields) interface{} {
	output := &model.List{}
	output.Limit = input.Limit
	output.Page = input.Page

	quantity, equipments, err := r.StockEquipmentService.List(input)
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
	err = json.Unmarshal(equipmentsByte, &output.StockEquipments)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, output)
}

func (r *resolver) ListByStockID(stockID string) interface{} {
	equipments, err := r.StockEquipmentService.ListByStockID(stockID)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, equipments)
}

func (r *resolver) GetByID(input *model.Field) interface{} {
	base, err := r.StockEquipmentService.GetByID(input)
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
	equipment, err := r.StockEquipmentService.GetByID(&model.Field{StockEquipmentID: input.StockEquipmentID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.StockEquipmentService.Update(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, equipment.StockEquipmentID)
}

func (r *resolver) Delete(input *model.Updated) interface{} {
	_, err := r.StockEquipmentService.GetByID(&model.Field{StockEquipmentID: input.StockEquipmentID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.StockEquipmentService.Delete(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, "Delete ok!")
}
