package repository_test

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"merchant_bank_payment_go_api/test/test_helpers"
	"os"
	"testing"
)

func CreateBlacklistTempFile() {
	fileContent, err := json.Marshal(test_helpers.ExpectedTokens)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(test_helpers.BlacklistTempFilename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func DeleteBlacklistTempFile() {
	err := os.Remove(test_helpers.BlacklistTempFilename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadBlackList_ShouldReturnBlackListToken(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	loadedToken, err := repo.LoadBlacklist()

	assert.Nil(t, err)
	assert.Equal(t, len(test_helpers.ExpectedTokens), len(loadedToken))
}

func TestLoadBlacklist_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, invalidFilename)

	tokenResults, err := repo.LoadBlacklist()

	assert.Nil(t, tokenResults)
	assert.NotNil(t, err)
}

func TestSaveBlacklist_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	newExpectedTokens := append(test_helpers.ExpectedTokens, "token4")

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	err := repo.SaveBlacklist(newExpectedTokens)

	fileContent, err := os.ReadFile(test_helpers.BlacklistTempFilename)
	assert.Nil(t, err)

	var tokenResults []string
	err = json.Unmarshal(fileContent, &tokenResults)
	assert.Nil(t, err)
	assert.Equal(t, newExpectedTokens, tokenResults)
}

func TestSaveBlacklist_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_blacklist.json"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, invalidFilename)

	err := repo.SaveBlacklist(test_helpers.ExpectedTokens)

	assert.NotNil(t, err)
}

func TestAddToBlacklist_ShouldAddNewToken(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	token := "token4"
	expectedBlacklistToken := []string{"token1", "token2", "token3", "token4"}

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	err := repo.AddToBlacklist(token)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(test_helpers.BlacklistTempFilename)
	assert.Nil(t, err)

	var tokenResults []string
	err = json.Unmarshal(fileContent, &tokenResults)
	assert.Nil(t, err)

	assert.Equal(t, expectedBlacklistToken, tokenResults)
}

func TestAddToBlacklist_ShouldReturnErrorWhenAlreadyBlacklisted(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	token := "token1"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	err := repo.AddToBlacklist(token)
	assert.NotNil(t, err)
}

func TestAddToBlacklist_ShouldReturnErrorIfLoadFails(t *testing.T) {
	invalidFilename := "nonexistent_folder/test_blacklist.json"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, invalidFilename)

	err := repo.AddToBlacklist("token4")

	assert.NotNil(t, err)
}

func TestIsTokenBlacklist_ShouldReturnTrue(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	token := "token1"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	blacklisted, err := repo.IsTokenBlacklisted(token)
	assert.True(t, blacklisted)
	assert.Nil(t, err)
}

func TestIsTokenBlacklist_ShouldReturnFalse(t *testing.T) {
	t.Cleanup(DeleteBlacklistTempFile)
	CreateBlacklistTempFile()

	token := "token4"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, test_helpers.BlacklistTempFilename)

	blacklisted, err := repo.IsTokenBlacklisted(token)
	assert.False(t, blacklisted)
	assert.Nil(t, err)
}
