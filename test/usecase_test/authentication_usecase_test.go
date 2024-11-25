package usecase_test

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
	"merchant_bank_payment_go_api/test/helper"
	"testing"
)

func TestLogin_ShouldReturnLoginResponse(t *testing.T) {
	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(helper.ExpectedCustomers[0], nil)

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: helper.ExpectedCustomers[0].Username,
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func TestLogin_ShouldReturnError_WhenInvalidUsername(t *testing.T) {
	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "susi").Return(entity.Customer{}, errors.New("customer not found"))

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "susi",
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogin_ShouldReturnError_WhenInvalidPassword(t *testing.T) {
	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(helper.ExpectedCustomers[0], nil)

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: helper.ExpectedCustomers[0].Username,
		Password: "password123",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogin_ShouldReturnError_WhenErrorLogOnInvalidUsername(t *testing.T) {
	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "susi").Return(entity.Customer{}, errors.New("customer not found"))

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file history not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "susi",
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogin_ShouldReturnError_WhenErrorLogOnInvalidPassword(t *testing.T) {
	customerId := uuid.New()
	mockCustomer := entity.Customer{
		Id:        customerId,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: helper.CreatedAt,
		UpdatedAt: helper.UpdatedAt,
	}

	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", "budi").Return(mockCustomer, nil)

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file history not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: "budi",
		Password: "password123",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogin_ShouldReturnError_WhenErrorLogOnFailedToLog(t *testing.T) {
	mockCustomerUseCase := new(helper.MockCustomerUseCase)
	mockCustomerUseCase.On("FindByUsername", helper.ExpectedCustomers[0].Username).Return(helper.ExpectedCustomers[0], nil)

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file history not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, mockCustomerUseCase, mockHistoryUseCase)

	request := model.LoginRequest{
		Username: helper.ExpectedCustomers[0].Username,
		Password: "password",
	}

	response, err := authUseCase.Login(request)

	assert.NotNil(t, err)
	assert.Empty(t, response.AccessToken)
}

func TestLogout_ShouldBlacklistToken(t *testing.T) {
	accessToken, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", accessToken).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", accessToken).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(accessToken)

	assert.Nil(t, err)
}

func TestLogout_ShouldReturnError_WhenAddToBlacklistFails(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "invalid_token").Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", "invalid_token").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout("invalid_token")

	assert.NotNil(t, err)
}

func TestLogout_ShouldReturnError_WhenErrorLog(t *testing.T) {
	accessToken, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", accessToken).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", accessToken).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(accessToken)

	assert.NotNil(t, err)
}

func TestLogout_ShouldReturnError_WhenErrorLogOnLogSuccessBlacklistToken(t *testing.T) {
	accessToken, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, "LOGOUT", "Customer ID extracted successfully", nil).Return(nil)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, "Token blacklisted successfully", nil).Return(errors.New("file not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", accessToken).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", accessToken).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(accessToken)

	assert.NotNil(t, err)
}

func TestLogout_ShouldReturnError_WhenErrorLogOnLogSuccessLogout(t *testing.T) {
	accessToken, _ := jwtutils.GenerateAccessToken(uuid.New().String())

	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, "LOGOUT", "Customer ID extracted successfully", nil).Return(nil)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, "Token blacklisted successfully", nil).Return(nil)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, "Logout successful", nil).Return(errors.New("file not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", accessToken).Return(false, nil)
	mockAuthRepository.On("AddToBlacklist", accessToken).Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.Logout(accessToken)

	assert.NotNil(t, err)
}

func TestLogout_ShouldReturnError_WhenErrorLogOnErrorExtractToken(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "accessToken").Return(true, errors.New("access token already blacklisted"))
	mockAuthRepository.On("AddToBlacklist", "accessToken").Return(errors.New("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)
	err := authUseCase.Logout("accessToken")

	assert.NotNil(t, err)
}

func TestLogout_ShouldReturnError_WhenErrorLogOnAddToBlacklistFails(t *testing.T) {
	accessToken, _ := jwtutils.GenerateAccessToken(uuid.New().String())
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, "LOGOUT", "Customer ID extracted successfully", nil).Return(nil)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("file not exists"))

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", accessToken).Return(true, errors.New("access token already blacklisted"))
	mockAuthRepository.On("AddToBlacklist", accessToken).Return(errors.New("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)
	err := authUseCase.Logout(accessToken)

	assert.NotNil(t, err)
}

func TestIsTokenBlacklisted_ShouldReturnTrue(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "blacklisted_token").Return(true, nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("blacklisted_token")

	assert.Nil(t, err)
	assert.True(t, isBlacklisted)
}

func TestIsTokenBlacklisted_ShouldReturnError_WhenRepositoryFails(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("IsTokenBlacklisted", "token_error").Return(false, fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	isBlacklisted, err := authUseCase.IsTokenBlacklisted("token_error")

	assert.NotNil(t, err)
	assert.False(t, isBlacklisted)
}

func TestAddToBlacklist_ShouldCallRepository(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "new_token").Return(nil)

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("new_token")

	assert.Nil(t, err)
	mockAuthRepository.AssertExpectations(t)
}

func TestAddToBlacklist_ShouldReturnError_WhenAddFails(t *testing.T) {
	mockHistoryUseCase := new(helper.MockHistoryUseCase)
	mockHistoryUseCase.On("LogAndAddHistory", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockAuthRepository := new(helper.MockAuthRepository)
	mockAuthRepository.On("AddToBlacklist", "token_error").Return(fmt.Errorf("repository error"))

	authUseCase := impl.NewAuthUseCaseImpl(mockAuthRepository, nil, mockHistoryUseCase)

	err := authUseCase.AddToBlacklist("token_error")

	assert.NotNil(t, err)
	mockAuthRepository.AssertExpectations(t)
}
