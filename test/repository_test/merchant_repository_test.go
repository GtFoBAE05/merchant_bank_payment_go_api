package repository_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"merchant_bank_payment_go_api/test/test_helpers"
	"os"
	"testing"
)

func CreateMerchantTempFile() {
	fileContent, err := json.Marshal(test_helpers.ExpectedMerchants)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(test_helpers.MerchantFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteMerchantTempFile() {
	err := os.Remove(test_helpers.MerchantFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadMerchant_ShouldReturnMerchantList(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	log := logrus.New()
	repo := impl.NewMerchantRepository(log, test_helpers.MerchantFilename)

	merchantResult, err := repo.LoadMerchants()

	assert.Nil(t, err)
	assert.Equal(t, len(test_helpers.ExpectedMerchants), len(merchantResult))
	assert.Equal(t, test_helpers.ExpectedMerchants[0].Id, merchantResult[0].Id)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].Name, merchantResult[0].Name)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].CreatedAt, merchantResult[0].CreatedAt)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].UpdatedAt, merchantResult[0].UpdatedAt)
}

func TestLoadMerchant_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewMerchantRepository(log, invalidFilename)

	merchantResult, err := repo.LoadMerchants()

	assert.Nil(t, merchantResult)
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	log := logrus.New()
	repo := impl.NewMerchantRepository(log, test_helpers.MerchantFilename)

	merchantResult, err := repo.FindById(test_helpers.MerchantId)
	assert.Nil(t, err)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].Id, merchantResult.Id)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].Name, merchantResult.Name)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].CreatedAt, merchantResult.CreatedAt)
	assert.Equal(t, test_helpers.ExpectedMerchants[0].UpdatedAt, merchantResult.UpdatedAt)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	merchantUuid := uuid.New()
	log := logrus.New()
	repo := impl.NewMerchantRepository(log, test_helpers.MerchantFilename)

	loadedMerchant, err := repo.FindById(merchantUuid)
	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, loadedMerchant)
}
