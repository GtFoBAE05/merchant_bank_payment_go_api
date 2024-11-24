package impl

import (
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
)

type PaymentTransactionImpl struct {
	Log      *logrus.Logger
	filename string
}

func NewPaymentTransactionImpl(log *logrus.Logger, filename string) *PaymentTransactionImpl {
	return &PaymentTransactionImpl{
		Log:      log,
		filename: filename,
	}
}

func (p PaymentTransactionImpl) LoadPayments() ([]entity.PaymentTransaction, error) {
	return nil, nil
}

func (p PaymentTransactionImpl) SavePayments(transactions []entity.PaymentTransaction) error {
	return nil
}

func (p PaymentTransactionImpl) AddPayment(payment entity.PaymentTransaction) error {
	return nil
}
