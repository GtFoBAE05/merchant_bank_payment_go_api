package repository_test

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

const customerFilename = "test_customers.json"

func CreateCustomerTempFile() {
	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		logrus.Error("Failed to parse UUID:", err)
		return
	}

	customers := []entity.Customer{{
		Id:        parsedUUID,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}}

	fileContent, err := json.Marshal(customers)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(customerFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteCustomerTempfile() {
	err := os.Remove(customerFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadCustomers_ShouldReturnCustomerList(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)

	CreateCustomerTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	expectedCustomers := []entity.Customer{{
		Id:        parsedUUID,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}}

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, customerFilename)

	loadedCustomers, err := repo.LoadCustomers()

	assert.Nil(t, err)
	assert.Equal(t, len(expectedCustomers), len(loadedCustomers))
	assert.Equal(t, expectedCustomers[0].Id, loadedCustomers[0].Id)
	assert.Equal(t, expectedCustomers[0].Username, loadedCustomers[0].Username)
	assert.Equal(t, expectedCustomers[0].Password, loadedCustomers[0].Password)
	assert.Equal(t, expectedCustomers[0].CreatedAt, loadedCustomers[0].CreatedAt)
	assert.Equal(t, expectedCustomers[0].UpdatedAt, loadedCustomers[0].UpdatedAt)
}

func TestLoadCustomers_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, invalidFilename)

	customers, err := repo.LoadCustomers()

	assert.Nil(t, customers)
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnCustomer(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)

	CreateCustomerTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	expectedCustomers := entity.Customer{
		Id:        parsedUUID,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, customerFilename)

	loadedCustomer, err := repo.FindById(parsedUUID)
	assert.Nil(t, err)
	assert.Equal(t, expectedCustomers.Id, loadedCustomer.Id)
	assert.Equal(t, expectedCustomers.Username, loadedCustomer.Username)
	assert.Equal(t, expectedCustomers.Password, loadedCustomer.Password)
	assert.Equal(t, expectedCustomers.CreatedAt, loadedCustomer.CreatedAt)
	assert.Equal(t, expectedCustomers.UpdatedAt, loadedCustomer.UpdatedAt)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)

	CreateCustomerTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-000000000000")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, customerFilename)

	loadedCustomer, err := repo.FindById(parsedUUID)
	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, loadedCustomer)
}

func TestFindByUsername_ShouldReturnCustomer(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)

	CreateCustomerTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	expectedCustomers := entity.Customer{
		Id:        parsedUUID,
		Username:  "budi",
		Password:  "$2a$10$2y2ss1Xs8TWZKWFS2//gnuhX/Ruhvx07lIN6jcZX1JziMvC/uLOJe",
		CreatedAt: "2024-11-22 11:31:58.769884426",
		UpdatedAt: "2024-11-22 11:31:58.769884426",
	}

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, customerFilename)

	loadedCustomer, err := repo.FindByUsername("budi")
	assert.Nil(t, err)
	assert.Equal(t, expectedCustomers.Id, loadedCustomer.Id)
	assert.Equal(t, expectedCustomers.Username, loadedCustomer.Username)
	assert.Equal(t, expectedCustomers.Password, loadedCustomer.Password)
	assert.Equal(t, expectedCustomers.CreatedAt, loadedCustomer.CreatedAt)
	assert.Equal(t, expectedCustomers.UpdatedAt, loadedCustomer.UpdatedAt)
}

func TestFindByUsername_ShouldReturnError(t *testing.T) {
	t.Cleanup(DeleteCustomerTempfile)

	CreateCustomerTempFile()

	log := logrus.New()
	repo := impl.NewCustomerRepository(log, customerFilename)

	loadedCustomer, err := repo.FindByUsername("aaa")
	assert.NotNil(t, err)
	assert.Equal(t, entity.Customer{}, loadedCustomer)
}
