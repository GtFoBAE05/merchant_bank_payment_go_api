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
)

func CreateCustomerTempFile() {
	fileContent, err := json.Marshal(helper.ExpectedCustomers)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(helper.CustomerFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteCustomerTempfile() {
	err := os.Remove(helper.CustomerFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadCustomers_ShouldReturnCustomerList(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)
	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResults, err := repo.LoadCustomers()

	assert.Nil(t, err)
	assert.Equal(t, len(helper.ExpectedCustomers), len(customerResults))
	assert.Equal(t, helper.ExpectedCustomers[0].Id, customerResults[0].Id)
	assert.Equal(t, helper.ExpectedCustomers[0].Username, customerResults[0].Username)
	assert.Equal(t, helper.ExpectedCustomers[0].Password, customerResults[0].Password)
	assert.Equal(t, helper.ExpectedCustomers[0].CreatedAt, customerResults[0].CreatedAt)
	assert.Equal(t, helper.ExpectedCustomers[0].UpdatedAt, customerResults[0].UpdatedAt)
}

func TestLoadCustomers_ShouldReturnError_WhenInvalidFile(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, invalidFilename)

	customerResults, err := repo.LoadCustomers()

	assert.Nil(t, customerResults)
	assert.NotNil(t, err)
}

func TestLoadCustomers_ShouldReturnError_WhenInvalidContent(t *testing.T) {
	err := os.WriteFile(helper.CustomerFilename, []byte(""), 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
		return
	}

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResults, err := repo.LoadCustomers()

	assert.Nil(t, customerResults)
	assert.NotNil(t, err)

	err = os.Remove(helper.CustomerFilename)
	if err != nil {
		logrus.Error("Error deleting to file:", err)
	}
}

func TestFindById_ShouldReturnCustomer(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)
	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResult, err := repo.FindById(helper.CustomerId)
	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedCustomers[0].Id, customerResult.Id)
	assert.Equal(t, helper.ExpectedCustomers[0].Username, customerResult.Username)
	assert.Equal(t, helper.ExpectedCustomers[0].Password, customerResult.Password)
	assert.Equal(t, helper.ExpectedCustomers[0].CreatedAt, customerResult.CreatedAt)
	assert.Equal(t, helper.ExpectedCustomers[0].UpdatedAt, customerResult.UpdatedAt)
}

func TestFindByCustomerId_ShouldReturnError_WhenNotFound(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)
	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResult, err := repo.FindById(uuid.New())
	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customerResult)
}

func TestFindByCustomerId_ShouldReturnError_WhenLoadCustomerError(t *testing.T) {
	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, "abc/invalidpath")

	_, err := repo.FindById(uuid.New())

	assert.NotNil(t, err)
}

func TestFindByUsername_ShouldReturnCustomer(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)
	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResult, err := repo.FindByUsername(helper.ExpectedCustomers[0].Username)
	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedCustomers[0].Id, customerResult.Id)
	assert.Equal(t, helper.ExpectedCustomers[0].Username, customerResult.Username)
	assert.Equal(t, helper.ExpectedCustomers[0].Password, customerResult.Password)
	assert.Equal(t, helper.ExpectedCustomers[0].CreatedAt, customerResult.CreatedAt)
	assert.Equal(t, helper.ExpectedCustomers[0].UpdatedAt, customerResult.UpdatedAt)
}

func TestFindByUsername_ShouldReturnError_WhenNotFound(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)
	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, helper.CustomerFilename)

	customerResult, err := repo.FindByUsername("aaa")
	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, customerResult)
}

func TestFindByUsername_ShouldReturnError_WhenLoadCustomerError(t *testing.T) {
	log := logrus.New()
	repo := impl.NewCustomerRepositoryImpl(log, "abc/invalidpath")

	_, err := repo.FindByUsername("budi")

	assert.NotNil(t, err)
}
