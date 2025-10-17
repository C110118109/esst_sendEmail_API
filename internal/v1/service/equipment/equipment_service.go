package equipment

import (
	"encoding/json"
	"esst_sendEmail/internal/pkg/log"
	"esst_sendEmail/internal/pkg/util"
	model "esst_sendEmail/internal/v1/structure/equipments"
)

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

	output.EquipmentID = util.GenerateUUID()
	output.CreatedTime = util.NowToUTC()

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

func (s *service) CreateBatch(input *model.BatchCreated) error {
	return s.Entity.CreateBatch(input)
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

func (s *service) ListByProjectID(projectID string) (output []*model.Base, err error) {
	fields, err := s.Entity.ListByProjectID(projectID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	marshal, err := json.Marshal(fields)
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
	field, err := s.Entity.GetByID(&model.Field{EquipmentID: input.EquipmentID})
	if err != nil {
		log.Error(err)

		return err
	}
	err = s.Entity.Update(field)

	return err
}

func (s *service) Update(input *model.Updated) (err error) {
	field, err := s.Entity.GetByID(&model.Field{EquipmentID: input.EquipmentID})
	if err != nil {
		log.Error(err)

		return err
	}

	marshal, err := json.Marshal(input)
	if err != nil {
		log.Error(err)

		return err
	}

	err = json.Unmarshal(marshal, &field)
	if err != nil {
		log.Error(err)

		return err
	}

	err = s.Entity.Update(field)

	return err
}
