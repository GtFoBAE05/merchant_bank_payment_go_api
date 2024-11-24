package repository

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"os"
	"testing"
)

const paymentTransactionTempFilename = "test_payment_transaction.json"

var fromCustomerId = uuid.New()
var toMerchantId = uuid.New()
var expectedPayments = []entity.PaymentTransaction{
	{
		Id:         uuid.New(),
		CustomerId: fromCustomerId,
		MerchantId: toMerchantId,
		Amount:     50000,
		Timestamp:  "2024-11-23 10:25:38.12345678",
	},
}

func CreatePaymentTransactionTempFile() {
	fileContent, err := json.Marshal(expectedPayments)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(paymentTransactionTempFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeletePaymentTransactionTempFile() {
	err := os.Remove(paymentTransactionTempFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadPaymentTransaction_ShouldReturnPaymentTransactionToken(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, paymentTransactionTempFilename)

	paymentTransactionResult, err := repo.LoadPayments()

	assert.Nil(t, err)
	assert.Equal(t, len(expectedPayments), len(paymentTransactionResult))
}

func TestLoadPaymentTransaction_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	paymentTransactionResult, err := repo.LoadPayments()

	assert.Nil(t, paymentTransactionResult)
	assert.NotNil(t, err)
}

func TestSavePaymentTransaction_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	newPayments := entity.PaymentTransaction{
		Id:         uuid.New(),
		CustomerId: fromCustomerId,
		MerchantId: toMerchantId,
		Amount:     10000,
		Timestamp:  "2024-11-24 10:25:38.12345678",
	}
	addedPayments := expectedPayments
	addedPayments = append(addedPayments, newPayments)

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, paymentTransactionTempFilename)

	err := repo.SavePayments(addedPayments)

	fileContent, err := os.ReadFile(paymentTransactionTempFilename)
	assert.Nil(t, err)

	var paymentTransactionsResult []entity.PaymentTransaction
	err = json.Unmarshal(fileContent, &paymentTransactionsResult)
	assert.Nil(t, err)
	assert.Equal(t, addedPayments, paymentTransactionsResult)
}

func TestSavePaymentTransaction_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_payment_transaction.json"

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	err := repo.SavePayments(expectedPayments)

	assert.NotNil(t, err)
}

func TestAddToPaymentTransaction_ShouldAddNewPayment(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	newPayments := entity.PaymentTransaction{
		Id:         uuid.New(),
		CustomerId: fromCustomerId,
		MerchantId: toMerchantId,
		Amount:     10000,
		Timestamp:  "2024-11-24 10:25:38.12345678",
	}
	addedPayments := expectedPayments
	addedPayments = append(addedPayments, newPayments)

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, paymentTransactionTempFilename)

	err := repo.AddPayment(newPayments)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(paymentTransactionTempFilename)
	assert.Nil(t, err)

	var paymentTransactionResult []entity.PaymentTransaction
	err = json.Unmarshal(fileContent, &paymentTransactionResult)
	assert.Nil(t, err)

	assert.Equal(t, addedPayments, paymentTransactionResult)
}

func TestAddToPaymentTransaction_ShouldReturnErrorWhenLoadFails(t *testing.T) {
	invalidFilename := "nonexistent_folder/test_payments_transaction.json"

	newPayments := entity.PaymentTransaction{
		Id:         uuid.New(),
		CustomerId: fromCustomerId,
		MerchantId: toMerchantId,
		Amount:     10000,
		Timestamp:  "2024-11-24 10:25:38.12345678",
	}

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	err := repo.AddPayment(newPayments)

	assert.NotNil(t, err)
}
