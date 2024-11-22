package repository

import (
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
)

type CustomerRepository interface {
	LoadCustomers() ([]entity.Customer, error)
	FindById(id uuid.UUID) (entity.Customer, error)
	FindByUsername(username string) (entity.Customer, error)
}
