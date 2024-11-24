package impl

import (
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type CustomerUseCaseImpl struct {
	HistoryUseCase     usecase.HistoryUseCase
	CustomerRepository repository.CustomerRepository
}

func NewCustomerUseCaseImpl(historyUseCase usecase.HistoryUseCase, customerRepository repository.CustomerRepository) *CustomerUseCaseImpl {
	return &CustomerUseCaseImpl{
		HistoryUseCase:     historyUseCase,
		CustomerRepository: customerRepository,
	}
}

func (c *CustomerUseCaseImpl) FindById(id string) (entity.Customer, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		logHistoryErr := c.handleLogHistory(id, "Failed to parse uuid", "Error parsing customer UUID", err)
		if logHistoryErr != nil {
			return entity.Customer{}, logHistoryErr
		}
		return entity.Customer{}, err
	}

	customer, err := c.CustomerRepository.FindById(parsedUUID)
	if err != nil {
		logHistoryErr := c.handleLogHistory(id, "Failed to find customer by id", err.Error(), err)
		if logHistoryErr != nil {
			return entity.Customer{}, logHistoryErr
		}
		return entity.Customer{}, err
	}

	logHistoryErr := c.handleLogHistory(id, "Successfully found customer by id", "Customer found successfully", nil)
	if logHistoryErr != nil {
		return entity.Customer{}, logHistoryErr
	}
	return customer, nil
}

func (c *CustomerUseCaseImpl) FindByUsername(username string) (entity.Customer, error) {
	customer, err := c.CustomerRepository.FindByUsername(username)
	if err != nil {
		logHistoryErr := c.handleLogHistory(username, "Failed to find customer by username", err.Error(), err)
		if logHistoryErr != nil {
			return entity.Customer{}, logHistoryErr
		}
		return entity.Customer{}, err
	}

	logHistoryErr := c.handleLogHistory(username, "Successfully found customer by username", "Customer found successfully", nil)
	if logHistoryErr != nil {
		return entity.Customer{}, logHistoryErr
	}
	return customer, nil
}

func (c *CustomerUseCaseImpl) handleLogHistory(idOrUsername, action, message string, err error) error {
	logHistoryErr := c.HistoryUseCase.LogAndAddHistory(idOrUsername, action, message, err)
	if logHistoryErr != nil {
		return logHistoryErr
	}
	return nil
}
