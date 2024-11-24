package repository

import (
	"merchant_bank_payment_go_api/internal/entity"
)

type PaymentTransactionRepository interface {
	LoadPayments() ([]entity.PaymentTransaction, error)
	SavePayments([]entity.PaymentTransaction) error
	AddPayment(payment entity.PaymentTransaction) error
}
