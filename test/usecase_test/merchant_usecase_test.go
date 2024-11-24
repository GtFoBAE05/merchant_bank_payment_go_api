package usecase

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

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	mockMerchantRepository := new(test_helpers.MockMerchantRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	merchantId := uuid.New()
	expectedMerchant := entity.Merchant{
		Id:        merchantId,
		Name:      "toko jaya",
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}

	mockMerchantRepository.On("FindById", merchantId).Return(expectedMerchant, nil)

	merchantResult, err := useCase.FindById(merchantId.String())

	assert.Nil(t, err)
	assert.Equal(t, expectedMerchant, merchantResult)
}

func TestFindById_ShouldReturnErrorParseToken(t *testing.T) {
	mockMerchantRepository := new(test_helpers.MockMerchantRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	merchantId := "nil"
	mockMerchantRepository.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	mockMerchantRepository := new(test_helpers.MockMerchantRepository)
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	useCase := impl.NewMerchantUseCaseImpl(mockHistoryUseCase, mockMerchantRepository)

	merchantId := uuid.New()
	mockMerchantRepository.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId.String())

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
}
