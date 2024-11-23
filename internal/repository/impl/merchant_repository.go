package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
)

type MerchantRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewMerchantRepository(log *logrus.Logger, filename string) *MerchantRepositoryImpl {
	return &MerchantRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (m *MerchantRepositoryImpl) LoadMerchant() ([]entity.Merchant, error) {
	return nil, nil
}

func (m *MerchantRepositoryImpl) FindById(id uuid.UUID) (entity.Merchant, error) {
	return entity.Merchant{}, nil
}
