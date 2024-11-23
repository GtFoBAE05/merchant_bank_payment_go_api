package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
	"time"
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
	h.Log.Infof("Attempting to add %s history to the histories", customerId)

	parsedCustomerUUID, err := uuid.Parse(customerId)
	if err != nil {
		h.Log.Errorf("Failed to parse customer uuid: %s", customerId)
		return err
	}

	newHistory := entity.History{
		Id:         uuid.New(),
		Action:     action,
		CustomerId: parsedCustomerUUID,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05.999999999"),
		Details:    details,
	}

	return h.HistoryRepository.AddHistory(newHistory)
}
