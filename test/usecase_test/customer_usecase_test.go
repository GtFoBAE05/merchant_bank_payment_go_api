package usecase_test

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"testing"
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

func TestFindById_ShouldReturnCustomer(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	log := logrus.New()
	useCase := impl.NewCustomerUseCase(log, mockRepo)

	customerId := uuid.New()
	expectedCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	mockRepo.On("FindById", customerId).Return(expectedCustomer, nil)

	customer, err := useCase.FindById(customerId.String())

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomer, customer)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	log := logrus.New()
	useCase := impl.NewCustomerUseCase(log, mockRepo)

	customerId := uuid.New()

	mockRepo.On("FindById", customerId).Return(entity.Customer{}, nil)

	customer, err := useCase.FindById(customerId.String())

	assert.Nil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}

func TestCustomerUseCase_FindByUsername(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	log := logrus.New()
	useCase := impl.NewCustomerUseCase(log, mockRepo)

	expectedCustomer := entity.Customer{
		Id:        uuid.New(),
		Username:  "budi",
		Password:  "hashedpassword",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	mockRepo.On("FindByUsername", "budi").Return(expectedCustomer, nil)

	customer, err := useCase.FindByUsername("budi")

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomer, customer)
}

func TestFindByUsername_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	log := logrus.New()
	useCase := impl.NewCustomerUseCase(log, mockRepo)

	username := "budi"

	mockRepo.On("FindByUsername", username).Return(entity.Customer{}, nil)

	customer, err := useCase.FindByUsername(username)

	assert.Nil(t, err)
	assert.Equal(t, entity.Customer{}, customer)
}
