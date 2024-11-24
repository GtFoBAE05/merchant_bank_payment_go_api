package repository

import (
	"merchant_bank_payment_go_api/internal/entity"
)

type PaymentTransactionRepository interface {
	LoadPayments() ([]entity.Payment, error)
	SavePayments([]entity.Payment) error
	AddPayment(payment entity.Payment) error
}
