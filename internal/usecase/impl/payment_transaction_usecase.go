package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type PaymentTransactionUseCaseImpl struct {
	Log                          *logrus.Logger
	PaymentTransactionRepository repository.PaymentTransactionRepository
	CustomerUseCase              usecase.CustomerUseCase
	MerchantUseCase              usecase.MerchantUseCase
}

func NewPaymentTransactionUseCaseImpl(log *logrus.Logger, transactionRepository repository.PaymentTransactionRepository,
	customerUseCase usecase.CustomerUseCase, merchantUseCase usecase.MerchantUseCase) *PaymentTransactionUseCaseImpl {
	return &PaymentTransactionUseCaseImpl{
		Log:                          log,
		PaymentTransactionRepository: transactionRepository,
		CustomerUseCase:              customerUseCase,
		MerchantUseCase:              merchantUseCase,
	}
}

func (p *PaymentTransactionUseCaseImpl) AddPayment(paymentRequest model.PaymentRequest) error {
	return nil
}
