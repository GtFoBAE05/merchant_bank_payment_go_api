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
	p.Log.Debugf("Loading payment transactions from file: %s", p.Filename)

	file, err := utils.ReadJsonFile(p.Filename, p.Log)
	if err != nil {
		p.Log.Errorf("Failed to read file %s: %v", p.Filename, err)
		return nil, fmt.Errorf("failed to read payment transactions file: %w", err)
	}

	var transactions []entity.Payment
	if err := json.Unmarshal(file, &transactions); err != nil {
		p.Log.Errorf("Failed to decode JSON from file %s: %v", p.Filename, err)
		return nil, fmt.Errorf("failed to parse payment transactions: %w", err)
	}

	p.Log.Infof("Successfully loaded %d payment transactions", len(transactions))
	return transactions, nil
}

func (p *PaymentTransactionImpl) SavePayments(transactions []entity.Payment) error {
	p.Log.Infof("Saving %d payment transactions to file: %s", len(transactions), p.Filename)

	if err := utils.WriteJSONFile(p.Filename, transactions, p.Log); err != nil {
		p.Log.Errorf("Error saving payment transactions to file %s: %v", p.Filename, err)
		return fmt.Errorf("failed to save payment transactions: %w", err)
	}

	p.Log.Infof("Successfully saved %d payment transactions", len(transactions))
	return nil
}

func (p *PaymentTransactionImpl) AddPayment(payment entity.Payment) error {
	transactions, err := p.LoadPayments()
	if err != nil {
		return err
	}

	p.Log.Infof("Adding new payment transaction with ID %s", payment.Id.String())
	transactions = append(transactions, payment)

	if err := p.SavePayments(transactions); err != nil {
		p.Log.Errorf("Failed to save updated payment transactions: %v", err)
		return fmt.Errorf("error saving updated payment transactions: %w", err)
	}

	p.Log.Infof("Payment transaction %s successfully added", payment.Id.String())
	return nil
}
