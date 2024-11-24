package usecase

import "merchant_bank_payment_go_api/internal/model"

type PaymentTransactionUseCase interface {
	AddPayment(customerId string, paymentRequest model.PaymentRequest) error
}
