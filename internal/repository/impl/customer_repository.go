package impl

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
)

type CustomerRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewCustomerRepository(log *logrus.Logger, filename string) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (r *CustomerRepositoryImpl) LoadCustomers() ([]entity.Customer, error) {
	return nil, nil
}

func (r *CustomerRepositoryImpl) FindById(id uuid.UUID) (entity.Customer, error) {
	return entity.Customer{}, errors.New("customer not found")
}

func (r *CustomerRepositoryImpl) FindByUsername(username string) (entity.Customer, error) {
	return entity.Customer{}, errors.New("customer not found")
}
