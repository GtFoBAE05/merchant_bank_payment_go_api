package repository_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"merchant_bank_payment_go_api/test/helper"
	"os"
	"testing"
	"time"
)

func CreatePaymentTransactionTempFile() {
	fileContent, err := json.Marshal(helper.ExpectedPayments)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(helper.PaymentTransactionTempFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeletePaymentTransactionTempFile() {
	err := os.Remove(helper.PaymentTransactionTempFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadPaymentTransaction_ShouldReturnPaymentTransactionToken(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, helper.PaymentTransactionTempFilename)

	paymentTransactionResult, err := repo.LoadPayments()

	assert.Nil(t, err)
	assert.Equal(t, len(helper.ExpectedPayments), len(paymentTransactionResult))
}

func TestLoadPaymentTransaction_ShouldReturnError_WhenInvalidFilename(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	paymentTransactionResult, err := repo.LoadPayments()

	assert.Nil(t, paymentTransactionResult)
	assert.NotNil(t, err)
}

func TestLoadPayment_ShouldReturnError_WhenInvalidContent(t *testing.T) {
	err := os.WriteFile(helper.PaymentTransactionTempFilename, []byte(""), 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
		return
	}

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, helper.PaymentTransactionTempFilename)

	paymentResult, err := repo.LoadPayments()

	assert.Nil(t, paymentResult)
	assert.NotNil(t, err)

	err = os.Remove(helper.PaymentTransactionTempFilename)
	if err != nil {
		logrus.Error("Error deleting to file:", err)
	}
}

func TestSavePaymentTransaction_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	newPayments := entity.Payment{
		Id:         uuid.New(),
		CustomerId: helper.CustomerId,
		MerchantId: helper.MerchantId,
		Amount:     10000,
		Timestamp:  helper.CreatedAt,
	}
	addedPayments := helper.ExpectedPayments
	addedPayments = append(addedPayments, newPayments)

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, helper.PaymentTransactionTempFilename)

	err := repo.SavePayments(addedPayments)

	fileContent, err := os.ReadFile(helper.PaymentTransactionTempFilename)
	assert.Nil(t, err)

	var paymentTransactionsResult []entity.Payment
	err = json.Unmarshal(fileContent, &paymentTransactionsResult)
	assert.Nil(t, err)
	assert.Equal(t, addedPayments, paymentTransactionsResult)
}

func TestSavePaymentTransaction_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_payment_transaction.json"

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	err := repo.SavePayments(helper.ExpectedPayments)

	assert.NotNil(t, err)
}

func TestAddToPaymentTransaction_ShouldAddNewPayment(t *testing.T) {
	t.Cleanup(DeletePaymentTransactionTempFile)
	CreatePaymentTransactionTempFile()

	newPayments := entity.Payment{
		Id:         uuid.New(),
		CustomerId: helper.CustomerId,
		MerchantId: helper.MerchantId,
		Amount:     10000,
		Timestamp:  helper.CreatedAt,
	}
	addedPayments := helper.ExpectedPayments
	addedPayments = append(addedPayments, newPayments)

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, helper.PaymentTransactionTempFilename)

	err := repo.AddPayment(newPayments)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(helper.PaymentTransactionTempFilename)
	assert.Nil(t, err)

	var paymentTransactionResult []entity.Payment
	err = json.Unmarshal(fileContent, &paymentTransactionResult)
	assert.Nil(t, err)

	assert.Equal(t, addedPayments, paymentTransactionResult)
}

func TestAddToPaymentTransaction_ShouldReturnError_WhenLoadFails(t *testing.T) {
	invalidFilename := "nonexistent_folder/test_payments_transaction.json"

	newPayments := entity.Payment{
		Id:         uuid.New(),
		CustomerId: helper.CustomerId,
		MerchantId: helper.MerchantId,
		Amount:     10000,
		Timestamp:  time.Now(),
	}

	log := logrus.New()
	repo := impl.NewPaymentTransactionImpl(log, invalidFilename)

	err := repo.AddPayment(newPayments)

	assert.NotNil(t, err)
}
