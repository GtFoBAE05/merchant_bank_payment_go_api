package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"testing"
)

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

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	mockRepo := new(MockMerchantRepository)
	log := logrus.New()
	useCase := impl.NewMerchantUseCaseImpl(log, mockRepo)

	merchantId := uuid.New()
	expectedMerchant := entity.Merchant{
		Id:        merchantId,
		Name:      "toko jaya",
		CreatedAt: "2024-11-22 12:00:00.769884426",
		UpdatedAt: "2024-11-22 12:00:00.769884426",
	}

	mockRepo.On("FindById", merchantId).Return(expectedMerchant, nil)

	merchantResult, err := useCase.FindById(merchantId.String())

	assert.Nil(t, err)
	assert.Equal(t, expectedMerchant, merchantResult)
}

func TestFindById_ShouldReturnErrorParseToken(t *testing.T) {
	mockRepo := new(MockMerchantRepository)
	log := logrus.New()
	useCase := impl.NewMerchantUseCaseImpl(log, mockRepo)

	merchantId := "nil"
	mockRepo.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId)

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockMerchantRepository)
	log := logrus.New()
	useCase := impl.NewMerchantUseCaseImpl(log, mockRepo)

	merchantId := uuid.New()
	mockRepo.On("FindById", merchantId).Return(entity.Merchant{}, errors.New("merchant not found"))

	merchantResult, err := useCase.FindById(merchantId.String())

	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, merchantResult)
}
