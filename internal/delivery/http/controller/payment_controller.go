package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	jwtutils "merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/usecase"
	"net/http"
)

type PaymentController struct {
	Log            *logrus.Logger
	PaymentUseCase usecase.PaymentTransactionUseCase
}

func NewPaymentController(log *logrus.Logger, paymentUseCase usecase.PaymentTransactionUseCase) *PaymentController {
	return &PaymentController{
		Log:            log,
		PaymentUseCase: paymentUseCase,
	}
}

func (p *PaymentController) AddPayment(c *gin.Context) {
	var paymentRequest model.PaymentRequest
	p.Log.Debug("Attempting add payment request")

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

	token, exists := c.Get("token")
	if !exists {
		p.Log.Warn("Token not found in context")
		c.JSON(http.StatusUnauthorized, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusUnauthorized,
			Message:    "Token not found",
			Data:       nil,
		})
		return
	}

	userId, err := jwtutils.ExtractIDFromToken(token.(string))
	if err != nil {
		return
	}

	err = p.PaymentUseCase.AddPayment(userId, paymentRequest)
	if err != nil {
		p.Log.Warnf("Error add payment: %s", err)
		c.JSON(http.StatusBadRequest, model.CommonResponse[interface{}]{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	p.Log.Infof("Successfully add payment")
	c.JSON(http.StatusOK, model.CommonResponse[interface{}]{
		HttpStatus: http.StatusOK,
		Message:    "Successfully add payment",
		Data:       nil,
	})
}
