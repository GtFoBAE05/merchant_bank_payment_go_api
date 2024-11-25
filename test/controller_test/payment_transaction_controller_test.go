package controller_test

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

func TestAddPayment_ShouldReturnSuccess(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully added payment",
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
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAddPayment_ShouldReturnError_WhenInvalidRequest(t *testing.T) {
	token, _ := utils.GenerateAccessToken(uuid.New().String())
	paymentRequest := model.PaymentRequest{
		Amount: 10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusBadRequest,
		Message:    "Invalid body request",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)

	log := logrus.New()
	paymentController := controller.NewPaymentTransactionController(log, mockPaymentTransactionUseCase)

	r := gin.Default()
	r.POST("/payment", paymentController.AddPayment)

	req := httptest.NewRequest("POST", "/payment", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}

func TestAddPayment_ShouldReturnError_WhenNotUserIdOnContext(t *testing.T) {
	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusUnauthorized,
		Message:    "User ID not found",
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

func TestAddPayment_ShouldReturnError_WhenInvalidMerchantId(t *testing.T) {
	customerId := uuid.New()
	token, _ := utils.GenerateAccessToken(customerId.String())
	merchantId := uuid.New()
	paymentRequest := model.PaymentRequest{
		MerchantId: merchantId.String(),
		Amount:     10000,
	}
	expectedCommonResponse := model.CommonResponse[interface{}]{
		HttpStatus: http.StatusBadRequest,
		Message:    "invalid merchant id",
		Data:       nil,
	}
	bodyJson, err := json.Marshal(paymentRequest)
	assert.Nil(t, err)

	mockPaymentTransactionUseCase := new(helper.MockPaymentTransactionUseCase)
	mockPaymentTransactionUseCase.On("AddPayment", customerId.String(), paymentRequest).Return(errors.New("invalid merchant id"))

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
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	response := new(model.CommonResponse[interface{}])
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, expectedCommonResponse.Message, response.Message)
	assert.Equal(t, expectedCommonResponse.HttpStatus, response.HttpStatus)
}
