package controller_test

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/delivery/http/controller"
	"merchant_bank_payment_go_api/internal/delivery/http/middleware"
	auth "merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Login(request model.LoginRequest) (model.LoginResponse, error) {
	args := m.Called(request)
	return args.Get(0).(model.LoginResponse), args.Error(1)
}

func (m *MockAuthUseCase) IsTokenBlacklisted(accessToken string) (bool, error) {
	args := m.Called(accessToken)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockAuthUseCase) AddToBlacklist(accessToken string) error {
	args := m.Called(accessToken)
	return args.Error(0)
}

func (m *MockAuthUseCase) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestLogin_ShouldReturnAccessToken(t *testing.T) {
	loginRequest := model.LoginRequest{
		Username: "budi",
		Password: "password",
	}
	loginResponse := model.LoginResponse{
		AccessToken: "accessToken",
	}
	commonResponse := model.CommonResponse[model.LoginResponse]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully logged in",
		Data:       loginResponse,
	}
	bodyJson, err := json.Marshal(loginRequest)
	assert.Nil(t, err)

	log := logrus.New()
	mockAuthUseCase := new(MockAuthUseCase)
	mockAuthUseCase.On("Login", loginRequest).Return(loginResponse, nil)

	authController := controller.NewAuthController(log, mockAuthUseCase)

	r := gin.Default()
	r.POST("/login", authController.Login)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := new(model.CommonResponse[model.LoginResponse])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, commonResponse.Message, response.Message)
	assert.Equal(t, commonResponse.HttpStatus, response.HttpStatus)
	assert.Equal(t, commonResponse.Data.AccessToken, response.Data.AccessToken)
}

func TestLogin_ShouldReturnError_WhenInvalidRequest(t *testing.T) {
	loginRequest := model.LoginRequest{
		Username: "budi",
	}
	commonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(loginRequest)
	assert.Nil(t, err)

	log := logrus.New()
	mockAuthUseCase := new(MockAuthUseCase)
	authController := controller.NewAuthController(log, mockAuthUseCase)

	r := gin.Default()
	r.POST("/login", authController.Login)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	response := new(model.CommonResponse[model.LoginResponse])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, commonResponse.Message, response.Message)
	assert.Equal(t, commonResponse.HttpStatus, response.HttpStatus)
}

func TestLogin_ShouldReturnError_WhenInvalidCredential(t *testing.T) {
	loginRequest := model.LoginRequest{
		Username: "budi",
		Password: "password",
	}
	commonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "invalid credential",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(loginRequest)
	assert.Nil(t, err)

	log := logrus.New()
	mockAuthUseCase := new(MockAuthUseCase)
	mockAuthUseCase.On("Login", loginRequest).Return(model.LoginResponse{}, errors.New("invalid credential"))

	authController := controller.NewAuthController(log, mockAuthUseCase)

	r := gin.Default()
	r.POST("/login", authController.Login)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, commonResponse.Message, response.Message)
	assert.Equal(t, commonResponse.HttpStatus, response.HttpStatus)
}

func TestLogout_ShouldReturnSuccess_WhenTokenIsValid(t *testing.T) {
	token, _ := auth.GenerateAccessToken(uuid.New().String())
	commonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully logged out",
		Data:       nil,
	}

	log := logrus.New()
	mockAuthUseCase := new(MockAuthUseCase)

	authController := controller.NewAuthController(log, mockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(false, nil)
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	r := gin.Default()
	r.Use(middleware.AuthMiddleware(mockAuthUseCase))
	r.POST("/logout", authController.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := new(model.CommonResponse[interface{}])
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, commonResponse.Message, response.Message)
	assert.Equal(t, commonResponse.HttpStatus, response.HttpStatus)
}

func TestLogout_ShouldReturnUnauthorized_WhenTokenNotFound(t *testing.T) {
	commonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "Token not found",
		Data:       nil,
	}
	log := logrus.New()
	mockAuthUseCase := new(MockAuthUseCase)
	authController := controller.NewAuthController(log, mockAuthUseCase)

	r := gin.Default()
	r.POST("/logout", authController.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	response := new(model.CommonResponse[interface{}])
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, commonResponse.HttpStatus, response.HttpStatus)
	assert.Equal(t, commonResponse.Message, response.Message)
}
