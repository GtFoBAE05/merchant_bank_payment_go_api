package usecase

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/test_helpers"
	"testing"
)

func TestLogin_ShouldReturnLoginResponse(t *testing.T) {
	customerId := uuid.New()
	mockCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}

	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "budi").Return(mockCustomer, nil)

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "budi",
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func TestLogin_ShouldReturnInvalidUsername(t *testing.T) {
	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "susi").Return(entity.Customer{}, errors.New("customer not found"))

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

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
		CreatedAt: test_helpers.CreatedAt,
		UpdatedAt: test_helpers.UpdatedAt,
	}

	mockCustomerUseCase := new(test_helpers.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "budi").Return(mockCustomer, nil)

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "budi",
		Password: "password123",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogout_ShouldBlacklistToken(t *testing.T) {
	newUuid, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", newUuid).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", newUuid).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(newUuid)

	assert.Nil(t, err)
}

func TestLogout_ShouldReturnError_WhenAddToBlacklistFails(t *testing.T) {
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "invalid_token").Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", "invalid_token").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout("invalid_token")

	assert.NotNil(t, err)
}

func TestIsTokenBlacklisted_ShouldReturnTrue(t *testing.T) {
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "blacklisted_token").Return(true, nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("blacklisted_token")

	assert.Nil(t, err)
	assert.True(t, isBlacklisted)
}

func TestIsTokenBlacklisted_ShouldReturnError_WhenRepositoryFails(t *testing.T) {
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "token_error").Return(false, fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("token_error")

	assert.NotNil(t, err)
	assert.False(t, isBlacklisted)
}

func TestAddToBlacklist_ShouldCallRepository(t *testing.T) {
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "new_token").Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("new_token")

	assert.Nil(t, err)
	mockAuthRepository.AssertExpectations(t)
}

func TestAddToBlacklist_ShouldReturnError_WhenAddFails(t *testing.T) {
	mockHistoryUseCase := new(test_helpers.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(test_helpers.MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "token_error").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("token_error")

	assert.NotNil(t, err)
	mockAuthRepository.AssertExpectations(t)
}
