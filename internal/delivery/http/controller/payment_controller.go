package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/usecase"
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

}
