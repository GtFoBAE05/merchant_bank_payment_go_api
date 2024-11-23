package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
)

type HistoryUseCaseImpl struct {
	Log               *logrus.Logger
	HistoryRepository repository.HistoryRepository
}

func NewHistoryUseCaseImpl(log *logrus.Logger, historyRepository repository.HistoryRepository) *HistoryUseCaseImpl {
	return &HistoryUseCaseImpl{
		Log:               log,
		HistoryRepository: historyRepository,
	}
}

func (h *HistoryUseCaseImpl) AddHistory(customerId, action, details string) error {
	return h.HistoryRepository.AddHistory(entity.History{})
}
