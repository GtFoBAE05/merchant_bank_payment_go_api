package usecase

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	jwtutils "merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"testing"
)

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

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) LoadBlacklist() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
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

func TestLogin_ShouldReturnLoginResponse(t *testing.T) {
	customerId := uuid.New()
	mockCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	log := logrus.New()

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "budi").Return(mockCustomer, nil)

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "budi",
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func TestLogin_ShouldReturnInvalidUsername(t *testing.T) {
	log := logrus.New()

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "susi").Return(entity.Customer{}, errors.New("customer not found"))

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "susi",
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogin_ShouldReturnInvalidPassword(t *testing.T) {
	customerId := uuid.New()
	mockCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	log := logrus.New()

	mockCustomerUseCase := new(MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "budi").Return(mockCustomer, nil)

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "budi",
		Password: "password123",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogout_ShouldBlacklistToken(t *testing.T) {
	u, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", u).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", u).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(u)

	assert.Nil(t, err)
}

func TestLogout_ShouldReturnError_WhenAddToBlacklistFails(t *testing.T) {
	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "invalid_token").Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", "invalid_token").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout("invalid_token")

	assert.NotNil(t, err)
}

func TestIsTokenBlacklisted_ShouldReturnTrue(t *testing.T) {
	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "blacklisted_token").Return(true, nil)

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("blacklisted_token")

	assert.Nil(t, err)
	assert.True(t, isBlacklisted)
}

func TestIsTokenBlacklisted_ShouldReturnError_WhenRepositoryFails(t *testing.T) {
	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "token_error").Return(false, fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("token_error")

	assert.NotNil(t, err)
	assert.False(t, isBlacklisted)
}

func TestAddToBlacklist_ShouldCallRepository(t *testing.T) {
	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "new_token").Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("new_token")

	assert.Nil(t, err)
	mockAuthRepository.AssertExpectations(t)
}

func TestAddToBlacklist_ShouldReturnError_WhenAddFails(t *testing.T) {
	log := logrus.New()

	mockHistoryUseCase := new(MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "token_error").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(log, mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("token_error")

	assert.NotNil(t, err)
	mockAuthRepository.AssertExpectations(t)
}
