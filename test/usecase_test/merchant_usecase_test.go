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

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	mockMerchantRepository.On("FindById", helper.MerchantId).Return(helper.ExpectedMerchants[0], nil)

	merchantResult, err := useCase.FindById(helper.MerchantId.String())

	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedMerchants[0], merchantResult)
	mockMerchantRepository.AssertExpectations(t)
}

func TestFindById_ShouldReturnErrorParseToken(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	merchantId := "nil"
	mockMerchantRepository.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
}

func TestFindById_ShouldReturnError_WhenMerchantNotFound(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	merchantId := uuid.New()
	mockMerchantRepository.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId.String())

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
	mockMerchantRepository.AssertExpectations(t)
}

func TestFindMerchantById_ShouldReturnError_WhenLogOnErrorParseToken(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	_, err := useCase.FindById("12345")
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnError_WhenLogOnReturnMerchantNotFound(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockMerchantRepository.On("FindById", helper.MerchantId).Return(entity.Merchant{}, errors.New("customer not found"))
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	_, err := useCase.FindById(helper.MerchantId.String())
	assert.NotNil(t, err)
	mockMerchantRepository.AssertExpectations(t)
}

func TestFindById_ShouldReturnError_WhenLogOnReturnMerchantError(t *testing.T) {
	mockMerchantRepository := new(helper.MockMerchantRepository)
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error On Log"))
	mockMerchantRepository.On("FindById", helper.MerchantId).Return(helper.ExpectedMerchants[0], nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	_, err := useCase.FindById(helper.MerchantId.String())
	assert.NotNil(t, err)
}
