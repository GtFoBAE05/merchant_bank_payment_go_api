package repository

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/repository/impl"
	"os"
	"testing"
)

const filename = "test_blacklist_token.json"

func createDataFile() {
	tokens := []string{"token1", "token2", "token3"}

	fileContent, err := json.Marshal(tokens)
	if err != nil {
		logrus.Error("Error marshalling data:", err)
		return
	}

	err = os.WriteFile(filename, fileContent, 0644)
	if err != nil {
		logrus.Error("Error writing to file:", err)
	}
}

func clearDataFile() {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("Error removing file:", err)
	}
}

func TestLoadBlackList_ShouldReturnBlackListToken(t *testing.T) {
	t.Cleanup(clearDataFile)
	createDataFile()

	expectedTokens := []string{"token1", "token2", "token3"}

	log := logrus.New()
	repo := impl.AuthRepositoryImpl{
		Filename: filename,
		Log:      log,
	}

	loadedToken, err := repo.LoadBlacklist()

	assert.Nil(t, err)
	assert.Equal(t, len(expectedTokens), len(loadedToken))
}

func TestLoadBlacklist_ShouldReturnError(t *testing.T) {
	invalidFilename := "empty.json"

	log := logrus.New()
	repo := impl.AuthRepositoryImpl{
		Filename: invalidFilename,
		Log:      log,
	}

	tokens, err := repo.LoadBlacklist()

	assert.Nil(t, tokens)
	assert.NotNil(t, err)
}

func TestSaveBlacklist_ShouldReturnSuccess(t *testing.T) {
	t.Cleanup(clearDataFile)
	createDataFile()

	expectedTokens := []string{"token1", "token2", "token3", "token4"}

	log := logrus.New()
	repo := impl.AuthRepositoryImpl{
		Log:      log,
		Filename: filename,
	}

	blacklistedTokens := []string{"token1", "token2", "token3", "token4"}
	err := repo.SaveBlacklist(blacklistedTokens)

	fileContent, err := os.ReadFile(filename)
	assert.Nil(t, err)

	var loadedTokens []string
	err = json.Unmarshal(fileContent, &loadedTokens)
	assert.Nil(t, err, "Error unmarshalling file content")
	assert.Equal(t, expectedTokens, loadedTokens)
}

func TestSaveBlacklist_ShouldReturnError(t *testing.T) {
	invalidFilename := "abc/test_blacklist.json"

	blacklistedTokens := []string{"token1", "token2", "token3"}

	log := logrus.New()
	repo := impl.NewAuthRepository(log, invalidFilename)

	err := repo.SaveBlacklist(blacklistedTokens)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error saving blacklist to file")
}

func TestAddToBlacklist_ShouldAddNewToken(t *testing.T) {
	t.Cleanup(clearDataFile)
	createDataFile()

	token := "token4"
	expectedBlacklistToken := []string{"token1", "token2", "token3", "token4"}

	log := logrus.New()
	repo := impl.NewAuthRepository(log, filename)

	err := repo.AddToBlacklist(token)
	assert.Nil(t, err)

	fileContent, err := os.ReadFile(filename)
	assert.Nil(t, err)

	var savedTokens []string
	err = json.Unmarshal(fileContent, &savedTokens)
	assert.Nil(t, err)

	assert.Equal(t, expectedBlacklistToken, savedTokens)
}

func TestAddToBlacklist_ShouldReturnErrorWhenAlreadyBlacklisted(t *testing.T) {
	t.Cleanup(clearDataFile)
	createDataFile()

	token := "token1"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, filename)

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
	t.Cleanup(clearDataFile)
	createDataFile()

	token := "token1"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, filename)

	blacklisted, err := repo.IsTokenBlacklisted(token)
	assert.True(t, blacklisted)
	assert.Nil(t, err)
}

func TestIsTokenBlacklist_ShouldReturnFalse(t *testing.T) {
	t.Cleanup(clearDataFile)
	createDataFile()

	token := "token4"

	log := logrus.New()
	repo := impl.NewAuthRepository(log, filename)

	blacklisted, err := repo.IsTokenBlacklisted(token)
	assert.False(t, blacklisted)
	assert.Nil(t, err)
}