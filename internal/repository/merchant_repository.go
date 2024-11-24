package repository

import (
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
)

type MerchantRepository interface {
	LoadMerchants() ([]entity.Merchant, error)
	FindById(id uuid.UUID) (entity.Merchant, error)
}
