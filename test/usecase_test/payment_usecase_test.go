package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"testing"
)

type MockPaymentTransactionRepository struct {
	mock.Mock
}

func (m *MockPaymentTransactionRepository) LoadPayments() ([]entity.PaymentTransaction, error) {
	args := m.Called()
	return args.Get(0).([]entity.PaymentTransaction), args.Error(1)
}

func (m *MockPaymentTransactionRepository) SavePayments(paymentTransactions []entity.PaymentTransaction) error {
	args := m.Called(paymentTransactions)
	return args.Error(0)
}

func (m *MockPaymentTransactionRepository) AddPayment(paymentTransaction entity.PaymentTransaction) error {
	args := m.Called(paymentTransaction)
	return args.Error(0)
}

type MockMerchantUseCase struct {
	mock.Mock
}

func (m *MockMerchantUseCase) FindById(id string) (entity.Merchant, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Merchant), args.Error(1)
}

func TestAddPayment_ShouldCallRepository(t *testing.T) {
	customerId := uuid.New()
	expectedCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}
	merchantId := uuid.New()
	expectedMerchant := entity.Merchant{
		Id:        merchantId,
		Name:      "toko jaya",
		CreatedAt: "2024-11-22 12:00:00.769884426",
		UpdatedAt: "2024-11-22 12:00:00.769884426",
	}

	mockPaymentRepository := new(MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId).Return(expectedCustomer)

	mockMerchantUseCase := new(MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", merchantId).Return(expectedMerchant)

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	log := logrus.New()
	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(log, mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase)

	err := paymentUseCase.AddPayment(paymentRequest)

	assert.Nil(t, err)
	mockPaymentRepository.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnErrorWhenInvalidCustomerId(t *testing.T) {
	customerId := uuid.New()
	merchantId := uuid.New()
	mockPaymentRepository := new(MockPaymentTransactionRepository)

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId).Return(errors.New("invalid customer"))

	mockMerchantUseCase := new(MockMerchantUseCase)

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	log := logrus.New()
	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(log, mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase)

	err := paymentUseCase.AddPayment(paymentRequest)

	assert.NotNil(t, err)
	mockPaymentRepository.AssertExpectations(t)
}

func TestAddPayment_ShouldReturnErrorWhenInvalidMerchantId(t *testing.T) {
	customerId := uuid.New()
	expectedCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}
	merchantId := uuid.New()

	mockPaymentRepository := new(MockPaymentTransactionRepository)
	mockPaymentRepository.On("AddPayment", mock.Anything).Return(nil)

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindById", customerId).Return(expectedCustomer)

	mockMerchantUseCase := new(MockMerchantUseCase)
	mockMerchantUseCase.On("FindById", merchantId).Return(errors.New("invalid merchant"))

	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}

	log := logrus.New()
	paymentUseCase := impl.NewPaymentTransactionUseCaseImpl(log, mockPaymentRepository, mockCustomerUseCase, mockMerchantUseCase)

	err := paymentUseCase.AddPayment(paymentRequest)

	assert.NotNil(t, err)
	mockPaymentRepository.AssertExpectations(t)
}
