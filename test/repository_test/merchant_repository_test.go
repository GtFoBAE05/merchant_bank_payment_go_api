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

const merchantFilename = "test_merchant.json"

func CreateMerchantTempFile() {
	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		logrus.Error("Failed to parse UUID:", err)
		return
	}

	merchant := []entity.Merchant{{
		Id:        parsedUUID,
		Name:      "toko jaya",
		CreatedAt: "2024-11-22 12:00:00.769884426",
		UpdatedAt: "2024-11-22 12:00:00.769884426",
	}}

	fileContent, err := json.Marshal(merchant)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(merchantFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteMerchantTempFile() {
	err := os.Remove(merchantFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadMerchant_ShouldReturnMerchantList(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	expectedMerchant := []entity.Merchant{{
		Id:        parsedUUID,
		Name:      "toko jaya",
		CreatedAt: "2024-11-22 12:00:00.769884426",
		UpdatedAt: "2024-11-22 12:00:00.769884426",
	}}

	log := logrus.New()
	repo := impl.MerchantRepositoryImpl{
		Filename: merchantFilename,
		Log:      log,
	}

	merchantResult, err := repo.LoadMerchant()

	assert.Nil(t, err)
	assert.Equal(t, len(expectedMerchant), len(merchantResult))
	assert.Equal(t, expectedMerchant[0].Id, merchantResult[0].Id)
	assert.Equal(t, expectedMerchant[0].Name, merchantResult[0].Name)
	assert.Equal(t, expectedMerchant[0].CreatedAt, merchantResult[0].CreatedAt)
	assert.Equal(t, expectedMerchant[0].UpdatedAt, merchantResult[0].UpdatedAt)
}

func TestLoadMerchant_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.MerchantRepositoryImpl{
		Filename: invalidFilename,
		Log:      log,
	}

	merchantResult, err := repo.LoadMerchant()

	assert.Nil(t, merchantResult)
	assert.NotNil(t, err)
}

func TestFindById_ShouldReturnMerchant(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	parsedUUID, err := uuid.Parse("685729de-cd87-4524-80bc-9b19cf58df22")
	if err != nil {
		t.Fatalf("Error parsing UUID: %v", err)
	}

	expectedMerchant := entity.Merchant{
		Id:        parsedUUID,
		Name:      "toko jaya",
		CreatedAt: "2024-11-22 12:00:00.769884426",
		UpdatedAt: "2024-11-22 12:00:00.769884426",
	}

	log := logrus.New()
	repo := impl.MerchantRepositoryImpl{
		Filename: merchantFilename,
		Log:      log,
	}

	merchantResult, err := repo.FindById(parsedUUID)
	assert.Nil(t, err)
	assert.Equal(t, expectedMerchant.Id, merchantResult.Id)
	assert.Equal(t, expectedMerchant.Name, merchantResult.Name)
	assert.Equal(t, expectedMerchant.CreatedAt, merchantResult.CreatedAt)
	assert.Equal(t, expectedMerchant.UpdatedAt, merchantResult.UpdatedAt)
}

func TestFindById_ShouldReturnError(t *testing.T) {
	t.Cleanup(DeleteMerchantTempFile)
	CreateMerchantTempFile()

	merchantUuid := uuid.New()

	log := logrus.New()
	repo := impl.MerchantRepositoryImpl{
		Filename: merchantFilename,
		Log:      log,
	}

	loadedMerchant, err := repo.FindById(merchantUuid)
	assert.NotNil(t, err)
	assert.Equal(t, entity.Merchant{}, loadedMerchant)
}
