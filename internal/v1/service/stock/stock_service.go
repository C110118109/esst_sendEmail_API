package stock

import (
	"encoding/json"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/stocks"
	"time"
)

// parseDate 解析日期字串，支援多種格式
func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// 支援的日期格式
	formats := []string{
		"2006-01-02",                // 標準日期格式
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"2006-01-02 15:04:05",       // 日期時間格式
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, nil // 解析失敗返回 nil，不報錯
}

func (s *service) Create(input *model.Created) (*model.Base, error) {
	var output model.Base

	marshal, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &output)
	if err != nil {
		return nil, err
	}

	output.StockID = util.GenerateUUID()
	output.CreatedTime = util.NowToUTC()

	// 解析日期
	if input.ExpectedDeliveryDate != "" {
		if parsedDate, err := parseDate(input.ExpectedDeliveryDate); err == nil && parsedDate != nil {
			output.ExpectedDeliveryDate = parsedDate
		}
	}
	if input.ContractStartDate != "" {
		if parsedDate, err := parseDate(input.ContractStartDate); err == nil && parsedDate != nil {
			output.ContractStartDate = parsedDate
		}
	}
	if input.ContractEndDate != "" {
		if parsedDate, err := parseDate(input.ContractEndDate); err == nil && parsedDate != nil {
			output.ContractEndDate = parsedDate
		}
	}

	table := &model.Table{}
	marshal, err = json.Marshal(output)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, table)
	if err != nil {
		return nil, err
	}

	err = s.Entity.Create(table)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (s *service) List(input *model.Fields) (quantity int64, output []*model.Base, err error) {
	amount, fields, err := s.Entity.List(input)
	if err != nil {
		log.Error(err)
		return 0, output, err
	}

	marshal, err := json.Marshal(fields)
	if err != nil {
		log.Error(err)
		return 0, output, err
	}

	err = json.Unmarshal(marshal, &output)
	if err != nil {
		log.Error(err)
		return 0, output, err
	}

	return amount, output, err
}

func (s *service) GetByID(input *model.Field) (output *model.Base, err error) {
	field, err := s.Entity.GetByID(input)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	marshal, err := json.Marshal(field)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = json.Unmarshal(marshal, &output)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return output, nil
}

func (s *service) Delete(input *model.Updated) (err error) {
	field, err := s.Entity.GetByID(&model.Field{StockID: input.StockID})
	if err != nil {
		log.Error(err)
		return err
	}
	err = s.Entity.Delete(&model.Field{StockID: field.StockID})
	return err
}

func (s *service) Update(input *model.Updated) (err error) {
	field, err := s.Entity.GetByID(&model.Field{StockID: input.StockID})
	if err != nil {
		log.Error(err)
		return err
	}

	// 處理基本資訊
	if input.StockName != "" {
		field.StockName = input.StockName
	}
	if input.ContactName != "" {
		field.ContactName = input.ContactName
	}
	if input.ContactPhone != "" {
		field.ContactPhone = input.ContactPhone
	}
	if input.ContactEmail != "" {
		field.ContactEmail = input.ContactEmail
	}
	if input.Owner != "" {
		field.Owner = input.Owner
	}
	if input.Remark != "" {
		field.Remark = input.Remark
	}

	// 處理交貨資訊
	if input.ExpectedDeliveryPeriod != "" {
		field.ExpectedDeliveryPeriod = input.ExpectedDeliveryPeriod
	}
	if input.ExpectedDeliveryDate != "" {
		if parsedDate, err := parseDate(input.ExpectedDeliveryDate); err == nil && parsedDate != nil {
			field.ExpectedDeliveryDate = parsedDate
		}
	}
	if input.ExpectedContractPeriod != "" {
		field.ExpectedContractPeriod = input.ExpectedContractPeriod
	}
	if input.ContractStartDate != "" {
		if parsedDate, err := parseDate(input.ContractStartDate); err == nil && parsedDate != nil {
			field.ContractStartDate = parsedDate
		}
	}
	if input.ContractEndDate != "" {
		if parsedDate, err := parseDate(input.ContractEndDate); err == nil && parsedDate != nil {
			field.ContractEndDate = parsedDate
		}
	}
	if input.DeliveryAddress != "" {
		field.DeliveryAddress = input.DeliveryAddress
	}
	if input.SpecialRequirements != "" {
		field.SpecialRequirements = input.SpecialRequirements
	}

	// 更新時間
	now := time.Now()
	field.UpdatedTime = &now

	err = s.Entity.Update(field)
	return err
}
