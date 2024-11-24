package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase"
	"net/http"
)

type PaymentTransactionController struct {
	Log            *logrus.Logger
	PaymentUseCase usecase.PaymentTransactionUseCase
}

func NewPaymentTransactionController(log *logrus.Logger, paymentUseCase usecase.PaymentTransactionUseCase) *PaymentTransactionController {
	return &PaymentTransactionController{
		Log:            log,
		PaymentUseCase: paymentUseCase,
	}
}

func (p *PaymentTransactionController) AddPayment(c *gin.Context) {
	var paymentRequest model.PaymentRequest
	p.Log.Debug("Attempting to add payment request")

	err := c.ShouldBind(&paymentRequest)
	if err != nil {
		p.Log.Errorf("Invalid payment body request: %v", err)
		c.JSON(http.StatusBadRequest, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       nil,
		})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		p.Log.Warn("User ID not found in context")
		c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusUnauthorized,
			Message:    "User ID not found",
			Data:       nil,
		})
		return
	}

	err = p.PaymentUseCase.AddPayment(userId.(string), paymentRequest)
	if err != nil {
		p.Log.Warnf("Error adding payment: %v", err)
		c.JSON(http.StatusBadRequest, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	p.Log.Infof("Successfully added payment")
	c.JSON(http.StatusOK, model.CommonResponse[interface{}]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully added payment",
		Data:       nil,
	})
}
