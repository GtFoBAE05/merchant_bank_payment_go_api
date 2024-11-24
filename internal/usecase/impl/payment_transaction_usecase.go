package impl

import (
	"fmt"
	"github.com/google/uuid"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
	"time"
)

type PaymentTransactionUseCaseImpl struct {
	PaymentTransactionRepository repository.PaymentTransactionRepository
	CustomerUseCase              usecase.CustomerUseCase
	MerchantUseCase              usecase.MerchantUseCase
	HistoryUseCase               usecase.HistoryUseCase
}

func NewPaymentTransactionUseCaseImpl(transactionRepository repository.PaymentTransactionRepository,
	customerUseCase usecase.CustomerUseCase, merchantUseCase usecase.MerchantUseCase, historyUseCase usecase.HistoryUseCase) *PaymentTransactionUseCaseImpl {
	return &PaymentTransactionUseCaseImpl{
		PaymentTransactionRepository: transactionRepository,
		CustomerUseCase:              customerUseCase,
		MerchantUseCase:              merchantUseCase,
		HistoryUseCase:               historyUseCase,
	}
}

func (p *PaymentTransactionUseCaseImpl) AddPayment(customerId string, paymentRequest model.PaymentRequest) error {
	customer, err := p.CustomerUseCase.FindById(customerId)
	if err != nil {
		return p.handleLogHistory("-", "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
	}

	merchant, err := p.MerchantUseCase.FindById(paymentRequest.MerchantId)
	if err != nil {
		return p.handleLogHistory(customer.Id.String(), "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
	}

	transaction := entity.Payment{
		Id:         uuid.New(),
		CustomerId: customer.Id,
		MerchantId: merchant.Id,
		Amount:     paymentRequest.Amount,
		Timestamp:  time.Now(),
	}

	err = p.PaymentTransactionRepository.AddPayment(transaction)
	if err != nil {
		return p.handleLogHistory(customer.Id.String(), "PAYMENT", fmt.Sprintf("Payment failed: %v", err), err)
	}

	return nil
}

func (p *PaymentTransactionUseCaseImpl) handleLogHistory(customerId, action, message string, err error) error {
	errLog := p.HistoryUseCase.LogAndAddHistory(customerId, action, message, err)
	if errLog != nil {
		return errLog
	}
	return err
}
