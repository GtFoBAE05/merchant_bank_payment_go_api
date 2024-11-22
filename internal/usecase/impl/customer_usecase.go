package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
)

type CustomerUseCase struct {
	Log                *logrus.Logger
	CustomerRepository repository.CustomerRepository
}

func NewCustomerUseCase(log *logrus.Logger, customerRepository repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		Log:                log,
		CustomerRepository: customerRepository,
	}
}

func (c *CustomerUseCase) FindById(id string) (entity.Customer, error) {
	return entity.Customer{}, nil
}

func (c *CustomerUseCase) FindByUsername(username string) (entity.Customer, error) {
	return entity.Customer{}, nil
}
