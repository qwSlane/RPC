package records

import (
	"rpc/internal/database"
	"rpc/internal/models"
	"rpc/internal/services/records/types"
)

//go:generate generator

//generator:gen
type records interface {
	SetNewRecord(*types.Record) error
	GetBestN(*types.BestLevelCount) (*models.Level, error)
}

type RecordsService struct {
	Storage storage.Storage
}

func NewRecordsService(storage storage.Storage) *RecordsService {

	return &RecordsService{
		Storage: storage,
	}
}

func (h *RecordsService) SetNewRecord(record *types.Record) error {
	err := h.Storage.SetNewRecord(record.Level, record.Username, record.Score)
	if err != nil {
		return err
	}
	return nil
}

func (h *RecordsService) GetBestN(data *types.BestLevelCount) (*models.Level, error) {
	users, err := h.Storage.GetBestN(int(data.Count), int(data.Level))
	if err != nil {
		return nil, err
	}

	message := &models.Level{
		Level:  data.Count,
		Scores: users,
	}

	return message, nil
}
