package usecase

import "merchant_bank_payment_go_api/internal/model"

type PaymentTransactionUseCase interface {
	AddPayment(paymentRequest model.PaymentRequest) error
}
