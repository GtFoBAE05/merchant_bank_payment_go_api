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

func CreateMerchantTempFile() {
	fileContent, err := json.Marshal(helper.ExpectedMerchants)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(helper.MerchantFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteMerchantTempFile() {
	err := os.Remove(helper.MerchantFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadMerchant_ShouldReturnMerchantList(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, helper.MerchantFilename)

	merchantResult, err := repo.LoadMerchants()

	assert.Nil(t, err)
	assert.Equal(t, len(helper.ExpectedMerchants), len(merchantResult))
	assert.Equal(t, helper.ExpectedMerchants[0].Id, merchantResult[0].Id)
	assert.Equal(t, helper.ExpectedMerchants[0].Name, merchantResult[0].Name)
	assert.Equal(t, helper.ExpectedMerchants[0].CreatedAt, merchantResult[0].CreatedAt)
	assert.Equal(t, helper.ExpectedMerchants[0].UpdatedAt, merchantResult[0].UpdatedAt)
}

func TestLoadMerchant_ShouldReturnError_WhenInvalidFileName(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, invalidFilename)

	merchantResult, err := repo.LoadMerchants()

	assert.Nil(t, merchantResult)
	assert.NotNil(t, err)
}

func TestLoadMerchant_ShouldReturnError_WhenInvalidContent(t *testing.T) {
	err := os.WriteFile(helper.MerchantFilename, []byte(""), 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
		return
	}

	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, helper.MerchantFilename)

	customerResults, err := repo.LoadMerchants()

	assert.Nil(t, customerResults)
	assert.NotNil(t, err)

	err = os.Remove(helper.MerchantFilename)
	if err != nil {
		logrus.Error("Error deleting to file:", err)
	}
}

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, helper.MerchantFilename)

	merchantResult, err := repo.FindById(helper.MerchantId)
	assert.Nil(t, err)
	assert.Equal(t, helper.ExpectedMerchants[0].Id, merchantResult.Id)
	assert.Equal(t, helper.ExpectedMerchants[0].Name, merchantResult.Name)
	assert.Equal(t, helper.ExpectedMerchants[0].CreatedAt, merchantResult.CreatedAt)
	assert.Equal(t, helper.ExpectedMerchants[0].UpdatedAt, merchantResult.UpdatedAt)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	merchantUuid := uuid.New()
	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, helper.MerchantFilename)

	loadedMerchant, err := repo.FindById(merchantUuid)
	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, loadedMerchant)
}

func TestFindById_ShouldReturnError_WhenLoadMerchantError(t *testing.T) {
	log := logrus.New()
	repo := impl.NewMerchantRepositoryImpl(log, "abc/invalidpath")

	_, err := repo.FindById(uuid.New())

	assert.NotNil(t, err)
}
