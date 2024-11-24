package impl

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
	"time"
)

type PaymentTransactionUseCaseImpl struct {
	Log                          *logrus.Logger
	PaymentTransactionRepository repository.PaymentTransactionRepository
	CustomerUseCase              usecase.CustomerUseCase
	MerchantUseCase              usecase.MerchantUseCase
	HistoryUseCase               usecase.HistoryUseCase
}

func NewPaymentTransactionUseCaseImpl(log *logrus.Logger, transactionRepository repository.PaymentTransactionRepository,
	customerUseCase usecase.CustomerUseCase, merchantUseCase usecase.MerchantUseCase, historyUseCase usecase.HistoryUseCase) *PaymentTransactionUseCaseImpl {
	return &PaymentTransactionUseCaseImpl{
		Log:                          log,
		PaymentTransactionRepository: transactionRepository,
		CustomerUseCase:              customerUseCase,
		MerchantUseCase:              merchantUseCase,
		HistoryUseCase:               historyUseCase,
	}
}

func (p *PaymentTransactionUseCaseImpl) AddPayment(customerId string, paymentRequest model.PaymentRequest) error {
	customer, err := p.CustomerUseCase.FindById(customerId)
	if err != nil {
		errLog := p.HistoryUseCase.LogAndAddHistory("-", "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
		if errLog != nil {
			return errLog
		}
		return err
	}

	merchant, err := p.MerchantUseCase.FindById(paymentRequest.MerchantId)
	if err != nil {
		errLog := p.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
		if errLog != nil {
			return errLog
		}
		return err
	}

	transaction := entity.PaymentTransaction{
		Id:         uuid.New(),
		CustomerId: customer.Id,
		MerchantId: merchant.Id,
		Amount:     paymentRequest.Amount,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05.999999999"),
	}

	err = p.PaymentTransactionRepository.AddPayment(transaction)
	if err != nil {
		errLog := p.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
		if errLog != nil {
			return errLog
		}
		return err
	}

	return nil
}
