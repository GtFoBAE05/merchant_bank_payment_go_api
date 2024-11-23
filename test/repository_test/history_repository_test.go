package repository

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"os"
	"testing"
)

const historyTempFilename = "test_history.json"

var customerId = uuid.New()
var merchantId = uuid.New()

var histories = []entity.History{
	{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 10:25:38.12345678",
		Details:    "Login successful",
	},
	{
		Id:         uuid.New(),
		Action:     "PAYMENT",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 14:48:22.987654321",
		Details:    fmt.Sprintf("Payment of 20000 to Merchant Id %s", merchantId),
	},
	{
		Id:         uuid.New(),
		Action:     "LOGOUT",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 19:02:55.654321987",
		Details:    "Logout successful",
	},
}

func CreateHistoryTempFile() {
	fileContent, err := json.Marshal(histories)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(historyTempFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteHistoryTempFile() {
	err := os.Remove(historyTempFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadHistories_ShouldReturnHistories(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, historyTempFilename)

	historiesResult, err := repo.LoadHistories()

	assert.Nil(t, err)
	assert.Equal(t, len(histories), len(historiesResult))
	assert.Equal(t, histories, historiesResult)
}

func TestLoadHistories_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	histories, err := repo.LoadHistories()

	assert.Nil(t, histories)
	assert.NotNil(t, err)
}

func TestSaveHistories_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, historyTempFilename)

	expectedHistories := histories
	expectedHistories = append(expectedHistories, entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 19:05:55.654321987",
		Details:    "Login successful",
	})
	err := repo.SaveHistories(expectedHistories)

	fileContent, err := os.ReadFile(historyTempFilename)
	assert.Nil(t, err)

	var historiesResult []entity.History
	err = json.Unmarshal(fileContent, &historiesResult)
	assert.Nil(t, err)
	assert.Equal(t, expectedHistories, historiesResult)
}

func TestSaveHistories_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_history.json"

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	err := repo.SaveHistories(histories)

	assert.NotNil(t, err)
}

func TestAddToHistory_ShouldAddNewHistory(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	newHistory := entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 19:05:55.654321987",
		Details:    "Login successful",
	}
	expectedHistories := histories
	expectedHistories = append(expectedHistories, newHistory)

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, historyTempFilename)

	err := repo.AddHistory(newHistory)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(historyTempFilename)
	assert.Nil(t, err)

	var historyResult []entity.History
	err = json.Unmarshal(fileContent, &historyResult)
	assert.Nil(t, err)

	assert.Equal(t, expectedHistories, historyResult)
}

func TestAddToHistory_ShouldReturnErrorWhenLoadFails(t *testing.T) {
	invalidFilename := "nonexistent_folder/test_history.json"

	newHistory := entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: customerId,
		Timestamp:  "2024-11-23 19:05:55.654321987",
		Details:    "Login successful",
	}

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	err := repo.AddHistory(newHistory)

	assert.NotNil(t, err)
}
