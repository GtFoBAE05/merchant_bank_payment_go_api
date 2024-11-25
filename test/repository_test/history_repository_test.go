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

func CreateHistoryTempFile() {
	fileContent, err := json.Marshal(helper.ExpectedHistories)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(helper.HistoryTempFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteHistoryTempFile() {
	err := os.Remove(helper.HistoryTempFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadHistories_ShouldReturnHistories(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, helper.HistoryTempFilename)

	historiesResult, err := repo.LoadHistories()

	assert.Nil(t, err)
	assert.Equal(t, len(helper.ExpectedHistories), len(historiesResult))
	assert.Equal(t, helper.ExpectedHistories, historiesResult)
}

func TestLoadHistories_ShouldReturnError_WhenInvalidFilename(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	historiesResult, err := repo.LoadHistories()

	assert.Nil(t, historiesResult)
	assert.NotNil(t, err)
}

func TestLoadHistories_ShouldReturnError_WhenInvalidContent(t *testing.T) {
	err := os.WriteFile(helper.HistoryTempFilename, []byte(""), 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
		return
	}

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, helper.HistoryTempFilename)

	historiesResult, err := repo.LoadHistories()

	assert.Nil(t, historiesResult)
	assert.NotNil(t, err)

	err = os.Remove(helper.HistoryTempFilename)
	if err != nil {
		logrus.Error("Error deleting to file:", err)
	}
}

func TestSaveHistories_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, helper.HistoryTempFilename)

	newExpectedHistories := helper.ExpectedHistories
	newExpectedHistories = append(newExpectedHistories, entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: uuid.New().String(),
		Timestamp:  helper.CreatedAt,
		Details:    "Login successful",
	})
	err := repo.SaveHistories(newExpectedHistories)

	fileContent, err := os.ReadFile(helper.HistoryTempFilename)
	assert.Nil(t, err)

	var historiesResult []entity.History
	err = json.Unmarshal(fileContent, &historiesResult)
	assert.Nil(t, err)
	assert.Equal(t, newExpectedHistories, historiesResult)
}

func TestSaveHistories_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_history.json"

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	err := repo.SaveHistories(helper.ExpectedHistories)

	assert.NotNil(t, err)
}

func TestAddToHistory_ShouldAddNewHistory(t *testing.T) {
	t.Cleanup(DeleteHistoryTempFile)
	CreateHistoryTempFile()

	newHistory := entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: uuid.New().String(),
		Timestamp:  helper.CreatedAt,
		Details:    "Login successful",
	}
	newExpectedHistories := helper.ExpectedHistories
	newExpectedHistories = append(newExpectedHistories, newHistory)

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, helper.HistoryTempFilename)

	err := repo.AddHistory(newHistory)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(helper.HistoryTempFilename)
	assert.Nil(t, err)

	var historyResult []entity.History
	err = json.Unmarshal(fileContent, &historyResult)
	assert.Nil(t, err)

	assert.Equal(t, newExpectedHistories, historyResult)
}

func TestAddToHistory_ShouldReturnErrorWhenLoadFails(t *testing.T) {
	invalidFilename := "nonexistent_folder/test_history.json"

	newHistory := entity.History{
		Id:         uuid.New(),
		Action:     "LOGIN",
		CustomerId: uuid.New().String(),
		Timestamp:  time.Now(),
		Details:    "Login successful",
	}

	log := logrus.New()
	repo := impl.NewHistoryRepositoryImpl(log, invalidFilename)

	err := repo.AddHistory(newHistory)

	assert.NotNil(t, err)
}
