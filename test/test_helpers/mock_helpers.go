package test_helpers

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) LoadCustomers() ([]entity.Customer, error) {
	args := m.Called()
	return args.Get(0).([]entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindById(id uuid.UUID) (entity.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByUsername(username string) (entity.Customer, error) {
	args := m.Called(username)
	return args.Get(0).(entity.Customer), args.Error(1)
}

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) LoadMerchants() ([]entity.Merchant, error) {
	args := m.Called()
	return args.Get(0).([]entity.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) FindById(id uuid.UUID) (entity.Merchant, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Merchant), args.Error(1)
}

type MockMerchantUseCase struct {
	mock.Mock
}

func (m *MockMerchantUseCase) FindById(id string) (entity.Merchant, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Merchant), args.Error(1)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) LoadBlacklist() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

type MockCustomerUseCase struct {
	mock.Mock
}

func (m *MockCustomerUseCase) FindById(id string) (entity.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Customer), args.Error(1)
}

func (m *MockCustomerUseCase) FindByUsername(username string) (entity.Customer, error) {
	args := m.Called(username)
	return args.Get(0).(entity.Customer), args.Error(1)
}

type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) LoadHistories() ([]entity.History, error) {
	args := m.Called()
	return args.Get(0).([]entity.History), args.Error(1)
}

func (m *MockHistoryRepository) SaveHistories(histories []entity.History) error {
	args := m.Called(histories)
	return args.Error(0)
}

func (m *MockHistoryRepository) AddHistory(history entity.History) error {
	args := m.Called(history)
	return args.Error(0)
}

type MockHistoryUseCase struct {
	mock.Mock
}

func (m *MockHistoryUseCase) AddHistory(customerId, action, details string) error {
	args := m.Called(customerId, action, details)
	return args.Error(0)
}

func (m *MockHistoryUseCase) LogAndAddHistory(userId, action, message string, err error) error {
	args := m.Called(userId, action, message, err)
	return args.Error(0)
}

func (m *MockAuthRepository) SaveBlacklist(blacklistedTokens []string) error {
	args := m.Called(blacklistedTokens)
	return args.Error(0)
}

func (m *MockAuthRepository) AddToBlacklist(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthRepository) IsTokenBlacklisted(token string) (bool, error) {
	args := m.Called(token)
	return args.Get(0).(bool), args.Error(1)
}

type MockPaymentTransactionRepository struct {
	mock.Mock
}

func (m *MockPaymentTransactionRepository) LoadPayments() ([]entity.Payment, error) {
	args := m.Called()
	return args.Get(0).([]entity.Payment), args.Error(1)
}

func (m *MockPaymentTransactionRepository) SavePayments(paymentTransactions []entity.Payment) error {
	args := m.Called(paymentTransactions)
	return args.Error(0)
}

func (m *MockPaymentTransactionRepository) AddPayment(paymentTransaction entity.Payment) error {
	args := m.Called(paymentTransaction)
	return args.Error(0)
}
