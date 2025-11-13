package stock

import (
	"encoding/json"
	"errors"
	"time"

	"esst_sendEmail/internal/pkg/code"
	"esst_sendEmail/internal/pkg/linebot"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/stocks"

	"gorm.io/gorm"
)

func (r *resolver) Create(trx *gorm.DB, input *model.Created) interface{} {
	defer trx.Rollback()

	stock, err := r.StockService.WithTrx(trx).Create(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	trx.Commit()

	// 建立成功後,發送 LINE 通知
	// 注意:現貨報備需要在設備建立完成後才發送通知
	// 所以實際通知會在 stock_equipment 的 CreateBatch 中發送

	return code.GetCodeMessage(code.Successful, stock.StockID)
}

func (r *resolver) List(input *model.Fields) interface{} {
	output := &model.List{}
	output.Limit = input.Limit
	output.Page = input.Page

	quantity, stocks, err := r.StockService.List(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	stocksByte, err := json.Marshal(stocks)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	output.Pages = util.Pagination(quantity, output.Limit)
	err = json.Unmarshal(stocksByte, &output.Stocks)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err.Error())
	}

	return code.GetCodeMessage(code.Successful, output)
}

func (r *resolver) GetByID(input *model.Field) interface{} {
	base, err := r.StockService.GetByID(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	frontStock := &model.Single{}
	stocksByte, _ := json.Marshal(base)
	err = json.Unmarshal(stocksByte, &frontStock)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, frontStock)
}

func (r *resolver) Update(input *model.Updated) interface{} {
	// 驗證現貨是否存在
	stock, err := r.StockService.GetByID(&model.Field{StockID: input.StockID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	// 執行更新
	err = r.StockService.Update(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, stock.StockID)
}

func (r *resolver) Delete(input *model.Updated) interface{} {
	// 驗證現貨是否存在
	_, err := r.StockService.GetByID(&model.Field{StockID: input.StockID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.GetCodeMessage(code.DoesNotExist, err)
		}

		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	err = r.StockService.Delete(input)
	if err != nil {
		log.Error(err)
		return code.GetCodeMessage(code.InternalServerError, err)
	}

	return code.GetCodeMessage(code.Successful, "Delete ok!")
}

// sendStockLineNotification 發送現貨報備 LINE 通知
func (r *resolver) sendStockLineNotification(stockID string) {
	log.Info("Preparing to send stock LINE notification for stock:", stockID)

	// 查詢現貨資訊
	stock, err := r.StockService.GetByID(&model.Field{StockID: stockID})
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
		StockID:                stockID,
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
