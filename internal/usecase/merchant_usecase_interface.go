package usecase

import "merchant_bank_payment_go_api/internal/entity"

type MerchantUseCase interface {
	FindById(id string) (entity.Merchant, error)
}
