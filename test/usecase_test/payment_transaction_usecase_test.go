package usecase_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/helper"
	"testing"
)

func TestAddPayment_ShouldCallRepository(t *testing.T) {
	mockPaymentRepository := new(helper.MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", helper.CustomerId.String()).Return(helper.ExpectedCustomers[0], nil)

	mockMerchantUseCase := new(helper.MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", helper.MerchantId.String()).Return(helper.ExpectedMerchants[0], nil)

	paymentRequest := model.PaymentRequest{
		MerchantId: helper.MerchantId.String(),
		Amount:     10000,
	}

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(helper.CustomerId.String(), paymentRequest)

	assert.Nil(t, err)
	mockPaymentRepository.AssertExpectations(t)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnError_WhenInvalidCustomerId(t *testing.T) {
	customerId := uuid.New()
	mockPaymentRepository := new(helper.MockPaymentTransactionRepository)

	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId.String()).Return(entity.Customer{}, errors.New("invalid customer"))

	mockMerchantUseCase := new(helper.MockMerchantUseCase)

	paymentRequest := model.PaymentRequest{
		MerchantId: helper.MerchantId.String(),
		Amount:     10000,
	}

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(customerId.String(), paymentRequest)

	assert.NotNil(t, err)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnError_WhenErrorLogOnFindMerchant(t *testing.T) {
	merchantId := uuid.New()

	mockPaymentRepository := new(helper.MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", helper.CustomerId.String()).Return(helper.ExpectedCustomers[0], nil)

	mockMerchantUseCase := new(helper.MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", merchantId.String()).Return(entity.Merchant{}, errors.New("merchant not found"))

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("something wrong"))

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(helper.CustomerId.String(), paymentRequest)

	assert.NotNil(t, err)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnError_WhenErrorLogOnAddPayment(t *testing.T) {
	mockPaymentRepository := new(helper.MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(errors.New("error add"))

	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", helper.CustomerId.String()).Return(helper.ExpectedCustomers[0], nil)

	mockMerchantUseCase := new(helper.MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", helper.MerchantId.String()).Return(helper.ExpectedMerchants[0], nil)

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("something wrong"))

	paymentRequest := model.PaymentRequest{
		MerchantId: helper.MerchantId.String(),
		Amount:     10000,
	}

	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase, mockHistoryUseCase)

	err := paymentUseCase.AddPayment(helper.CustomerId.String(), paymentRequest)

	assert.NotNil(t, err)
	mockCustomerUseCase.AssertExpectations(t)
	mockMerchantUseCase.AssertExpectations(t)
}
