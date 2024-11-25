package usecase_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/helper"
	"testing"
)

func TestFindById_ShouldReturnCustomer(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	mockCustomerRepository.On("FindById", helper.CustomerId).Return(helper.ExpectedCustomers[0], nil)

	customer, err := useCase.FindById(helper.CustomerId.String())

	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedCustomers[0], customer)
	mockCustomerRepository.AssertExpectations(t)
}

func TestFindById_ShouldReturnError_WhenFailedParseToken(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	customerId := "abcdef"

	mockCustomerRepository.On("FindById", customerId).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindById(customerId)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestFindById_ShouldReturnError_WhenCustomerNotFound(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	customerId := uuid.New()

	mockCustomerRepository.On("FindById", customerId).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindById(customerId.String())

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestFindById_ShouldReturnError_WhenLogOnErrorParseToken(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	_, err := useCase.FindById("12345")
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnError_WhenLogOnReturnCustomerNotFound(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockCustomerRepository.On("FindById", helper.CustomerId).Return(entity.Customer{}, errors.New("customer not found"))
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	_, err := useCase.FindById(helper.CustomerId.String())
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnError_WhenLogOnReturnCustomerError(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockCustomerRepository.On("FindById", helper.CustomerId).Return(helper.ExpectedCustomers[0], nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	_, err := useCase.FindById(helper.CustomerId.String())
	assert.NotNil(t, err)
}

func TestFindByUsername_ShouldReturnCustomer(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	mockCustomerRepository.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(helper.ExpectedCustomers[0], nil)

	customer, err := useCase.FindByUsername(helper.ExpectedCustomers[0].Username)

	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedCustomers[0], customer)
	mockCustomerRepository.AssertExpectations(t)
}

func TestFindByUsername_ShouldReturnError_WhenCustomerNotFound(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	username := "budi"

	mockCustomerRepository.On("FindByUsername", username).Return(entity.Customer{}, errors.New("customer not found"))

	customer, err := useCase.FindByUsername(username)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestFindByUsername_ShouldReturnError_WhenLogOnReturnCustomerNotFound(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockCustomerRepository.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(entity.Customer{}, errors.New("customer not found"))
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	_, err := useCase.FindByUsername(helper.ExpectedCustomers[0].Username)
	assert.NotNil(t, err)
}

func TestFindByUsername_ShouldReturnError_WhenLogOnReturnCustomerError(t *testing.T) {
	mockCustomerRepository := new(helper.MockCustomerRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockCustomerRepository.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(helper.ExpectedCustomers[0], nil)
	useCase := impl.NewCustomerUseCaseImpl(mockHistoryUseCase, mockCustomerRepository)

	_, err := useCase.FindByUsername(helper.ExpectedCustomers[0].Username)
	assert.NotNil(t, err)
}
