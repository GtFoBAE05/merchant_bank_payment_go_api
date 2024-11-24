package usecase_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/test_helpers"
	"testing"
)

func TestFindById_ShouldReturnCustomer(t *testing.T) {
	mockCustomerRepository := new(test_helpers.MockCustomerRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	mockCustomerRepository.On("FindById", test_helpers.CustomerId).Return(test_helpers.ExpectedCustomers[0], nil)

	customer, err := useCase.FindById(test_helpers.CustomerId.String())

	assert.Nil(t, err)
	assert.Equal(t, test_helpers.ExpectedCustomers[0], customer)
}

func TestFindById_ShouldReturnErrorParseToken(t *testing.T) {
	mockCustomerRepository := new(test_helpers.MockCustomerRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	customerId := "abcdef"

	mockCustomerRepository.On("FindById", customerId).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindById(customerId)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	mockCustomerRepository := new(test_helpers.MockCustomerRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	customerId := uuid.New()

	mockCustomerRepository.On("FindById", customerId).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindById(customerId.String())

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestCustomerUseCase_FindByUsername(t *testing.T) {
	mockCustomerRepository := new(test_helpers.MockCustomerRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	mockCustomerRepository.On("FindByUsername", test_helpers.ExpectedCustomers[0].Username).Return(test_helpers.ExpectedCustomers[0], nil)

	customer, err := useCase.FindByUsername(test_helpers.ExpectedCustomers[0].Username)

	assert.Nil(t, err)
	assert.Equal(t, test_helpers.ExpectedCustomers[0], customer)
}

func TestFindByUsername_ShouldReturnError(t *testing.T) {
	mockCustomerRepository := new(test_helpers.MockCustomerRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	username := "budi"

	mockCustomerRepository.On("FindByUsername", username).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindByUsername(username)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}
