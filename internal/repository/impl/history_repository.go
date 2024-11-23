package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
)

type HistoryRepositoryImpl struct {
	Log      *logrus.Logger
	filename string
}

func NewHistoryRepositoryImpl(log *logrus.Logger, filename string) *HistoryRepositoryImpl {
	return &HistoryRepositoryImpl{
		Log:      log,
		filename: filename,
	}
}

func (h *HistoryRepositoryImpl) LoadHistories() ([]entity.History, error) {
	return nil, nil
}

func (h *HistoryRepositoryImpl) SaveHistories(histories []entity.History) error {
	return nil

}

func (h *HistoryRepositoryImpl) AddHistory(history entity.History) error {
	return nil
}
