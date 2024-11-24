package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/test_helpers"
	"testing"
)

func TestAddPayment_ShouldCallRepository(t *testing.T) {
	customerId := uuid.New()
	expectedCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}
	merchantId := uuid.New()
	expectedMerchant := entity.Merchant{
		Id:        merchantId,
		Name:      "toko jaya",
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}

	mockPaymentRepository := new(test_helpers.MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId.String()).Return(expectedCustomer, nil)

	mockMerchantUseCase := new(test_helpers.MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", merchantId.String()).Return(expectedMerchant, nil)

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(customerId.String(), paymentRequest)

	assert.Nil(t, err)
	mockPaymentRepository.AssertExpectations(t)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnErrorWhenInvalidCustomerId(t *testing.T) {
	customerId := uuid.New()
	merchantId := uuid.New()
	mockPaymentRepository := new(test_helpers.MockPaymentTransactionRepository)

	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId.String()).Return(entity.Customer{}, errors.New("invalid customer"))

	mockMerchantUseCase := new(test_helpers.MockMerchantUseCase)

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(customerId.String(), paymentRequest)

	assert.NotNil(t, err)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnErrorWhenInvalidMerchantId(t *testing.T) {
	customerId := uuid.New()
	expectedCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}
	merchantId := uuid.New()

	mockPaymentRepository := new(test_helpers.MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId.String()).Return(expectedCustomer, nil)

	mockMerchantUseCase := new(test_helpers.MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", merchantId.String()).Return(entity.Merchant{}, errors.New("invalid merchant"))

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(customerId.String(), paymentRequest)

	assert.NotNil(t, err)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}
