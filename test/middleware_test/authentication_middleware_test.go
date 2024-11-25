package middleware_test

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/delivery/http/controller"
	"merchant_bank_payment_go_api/internal/delivery/http/middleware"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/utils"
	"merchant_bank_payment_go_api/test/helper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticationMiddleware_ShouldReturnError_WhenNoHeader(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "Authorization Header is required",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(nil)

	mockAuthUseCase := new(helper.MockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(false, nil)
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.Use(middleware.AuthenticationMiddleware(mockAuthUseCase))
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAuthenticationMiddleware_ShouldReturnError_WhenInvalidHeaderFormat(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "Invalid Authorization Header format, must be 'Bearer <token>'",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(nil)

	mockAuthUseCase := new(helper.MockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(false, nil)
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.Use(middleware.AuthenticationMiddleware(mockAuthUseCase))
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", " ")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAuthenticationMiddleware_ShouldReturnError_WhenHaveHeaderButNoToken(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "Invalid or expired token",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(nil)

	mockAuthUseCase := new(helper.MockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(false, nil)
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.Use(middleware.AuthenticationMiddleware(mockAuthUseCase))
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer ")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAuthenticationMiddleware_ShouldReturnError_WhenErrorCheckIsTokenBlacklisted(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusInternalServerError,
		Message:    "Internal server error",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(nil)

	mockAuthUseCase := new(helper.MockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(false, errors.New("internal error"))
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.Use(middleware.AuthenticationMiddleware(mockAuthUseCase))
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAuthenticationMiddleware_ShouldReturnError_WhenTokenAlreadyBlacklisted(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusForbidden,
		Message:    "Token is already blacklisted",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(nil)

	mockAuthUseCase := new(helper.MockAuthUseCase)
	mockAuthUseCase.On("IsTokenBlacklisted", token).Return(true, nil)
	mockAuthUseCase.On("AddToBlacklist", token).Return(nil)
	mockAuthUseCase.On("Logout", token).Return(nil)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.Use(middleware.AuthenticationMiddleware(mockAuthUseCase))
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}
