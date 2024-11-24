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
		Timestamp:  time.Now(),
		Details:    details,
	}

	return h.HistoryRepository.AddHistory(newHistory)
}

func (h *HistoryUseCaseImpl) LogAndAddHistory(userId, action, message string, err error) error {
	if err != nil {
		h.Log.Errorf(message+": %v", err)
	} else {
		h.Log.Infof(message)
	}

	historyErr := h.AddHistory(userId, action, message)
	if historyErr != nil {
		h.Log.Errorf("Failed to add history for %s: %v", action, historyErr)
	}
	return historyErr
}
