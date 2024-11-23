package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
)

type MerchantUseCaseImpl struct {
	Log                *logrus.Logger
	MerchantRepository repository.MerchantRepository
}

func NewMerchantUseCaseImpl(log *logrus.Logger, merchantRepository repository.MerchantRepository) *MerchantUseCaseImpl {
	return &MerchantUseCaseImpl{
		Log:                log,
		MerchantRepository: merchantRepository,
	}
}

func (m *MerchantUseCaseImpl) FindById(id string) (entity.Merchant, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		m.Log.Errorf("Failed to parse uuid: %s", id)
		return entity.Merchant{}, err
	}
	merchant, err := m.MerchantRepository.FindById(parsedUUID)
	if err != nil {
		m.Log.Errorf("Failed to find merchant with id %s: %v", id, err)
		return entity.Merchant{}, err
	}

	return merchant, nil
}
