package impl

import (
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type MerchantUseCaseImpl struct {
	HistoryUseCase     usecase.HistoryUseCase
	MerchantRepository repository.MerchantRepository
}

func NewMerchantUseCaseImpl(historyUseCase usecase.HistoryUseCase, merchantRepository repository.MerchantRepository) *MerchantUseCaseImpl {
	return &MerchantUseCaseImpl{
		HistoryUseCase:     historyUseCase,
		MerchantRepository: merchantRepository,
	}
}

func (m *MerchantUseCaseImpl) FindById(id string) (entity.Merchant, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		logHistoryErr := m.HistoryUseCase.LogAndAddHistory(id, "Failed to parse UUID", "Error parsing merchant UUID", err)
		if logHistoryErr != nil {
			return entity.Merchant{}, logHistoryErr
		}
		return entity.Merchant{}, err
	}

	merchant, err := m.MerchantRepository.FindById(parsedUUID)
	if err != nil {
		logHistoryErr := m.HistoryUseCase.LogAndAddHistory(id, "Failed to find merchant by id", err.Error(), err)
		if logHistoryErr != nil {
			return entity.Merchant{}, logHistoryErr
		}
		return entity.Merchant{}, err
	}

	logHistoryErr := m.HistoryUseCase.LogAndAddHistory(id, "Successfully found merchant by id", "Merchant found successfully", nil)
	if logHistoryErr != nil {
		return entity.Merchant{}, logHistoryErr
	}
	return merchant, nil
}
