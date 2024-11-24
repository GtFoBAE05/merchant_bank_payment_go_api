package impl

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/utils"
)

type PaymentTransactionImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewPaymentTransactionImpl(log *logrus.Logger, filename string) *PaymentTransactionImpl {
	return &PaymentTransactionImpl{
		Log:      log,
		Filename: filename,
	}
}

func (p *PaymentTransactionImpl) LoadPayments() ([]entity.Payment, error) {
	p.Log.Debugf("Loading payment transaction from file: %s", p.Filename)

	file, err := utils.ReadJsonFile(p.Filename, p.Log)
	if err != nil {
		p.Log.Errorf("Error reading file %s: %v", p.Filename, err)
		return nil, err
	}

	p.Log.Tracef("File content: %s", string(file))

	var paymentTransaction []entity.Payment
	err = json.Unmarshal(file, &paymentTransaction)
	if err != nil {
		p.Log.Errorf("Error decoding JSON from file %s: %v", p.Filename, err)
		return nil, err
	}

	p.Log.Infof("Successfully loaded %d payment transaction from %s", len(paymentTransaction), p.Filename)
	return paymentTransaction, nil
}

func (p *PaymentTransactionImpl) SavePayments(transactions []entity.Payment) error {
	p.Log.Infof("Saving %d payment transaction to file: %s", len(transactions), p.Filename)

	err := utils.WriteJSONFile(p.Filename, transactions, p.Log)
	if err != nil {
		p.Log.Errorf("Error saving payment transactions to file %s: %v", p.Filename, err)
		return fmt.Errorf("error saving payment transactions to file %s: %v", p.Filename, err)
	}

	p.Log.Infof("Successfully saved payment transactions to %s", p.Filename)
	return nil
}

func (p *PaymentTransactionImpl) AddPayment(payment entity.Payment) error {
	paymentTransactions, err := p.LoadPayments()
	if err != nil {
		return err
	}

	p.Log.Infof("Adding payment %s to payment transaction", payment)

	paymentTransactions = append(paymentTransactions, payment)

	return p.SavePayments(paymentTransactions)
}
