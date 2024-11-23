package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
)

type MerchantUseCaseImpl struct {
	log                *logrus.Logger
	MerchantRepository repository.MerchantRepository
}

func NewMerchantUseCaseImpl(log *logrus.Logger, merchantRepository repository.MerchantRepository) *MerchantUseCaseImpl {
	return &MerchantUseCaseImpl{
		log:                log,
		MerchantRepository: merchantRepository,
	}
}

func (m MerchantUseCaseImpl) FindById(id string) (entity.Merchant, error) {
	return entity.Merchant{}, nil
}
