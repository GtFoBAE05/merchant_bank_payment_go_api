package usecase

import "merchant_bank_payment_go_api/internal/entity"

type CustomerUseCaseInterface interface {
	FindById(id string) (entity.Customer, error)
	FindByUsername(username string) (entity.Customer, error)
}
