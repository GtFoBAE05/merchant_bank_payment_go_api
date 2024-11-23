package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
)

type CustomerUseCaseImpl struct {
	Log                *logrus.Logger
	CustomerRepository repository.CustomerRepository
}

func NewCustomerUseCaseImpl(log *logrus.Logger, customerRepository repository.CustomerRepository) *CustomerUseCaseImpl {
	return &CustomerUseCaseImpl{
		Log:                log,
		CustomerRepository: customerRepository,
	}
}

func (c *CustomerUseCaseImpl) FindById(id string) (entity.Customer, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		c.Log.Errorf("Failed to parse uuid: %s", id)
		return entity.Customer{}, err
	}
	customer, err := c.CustomerRepository.FindById(parsedUUID)
	if err != nil {
		c.Log.Errorf("Failed to find customer with id %s: %v", id, err)
		return entity.Customer{}, err
	}

	return customer, nil
}

func (c *CustomerUseCaseImpl) FindByUsername(username string) (entity.Customer, error) {
	c.Log.Debugf("Finding customer by username: %s", username)
	customer, err := c.CustomerRepository.FindByUsername(username)
	if err != nil {
		c.Log.Errorf("Failed to find customer with username %s: %v", username, err)
		return entity.Customer{}, err
	}

	c.Log.Infof("Successfully found customer with username: %s", username)
	return customer, nil
}
