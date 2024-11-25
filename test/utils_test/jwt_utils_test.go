package utils_test

import (
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/utils"
	"merchant_bank_payment_go_api/test/helper"
	"testing"
)

func TestGenerateAccessToken_ShouldReturnAccessToken(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	token, err := utils.GenerateAccessToken(helper.CustomerId.String())
	assert.Nil(t, err)

	idFromToken, err := utils.ExtractIDFromToken(token)
	assert.Nil(t, err)

	assert.Equal(t, helper.CustomerId.String(), idFromToken)
}

func TestVerifyAccessToken_ShouldReturnTrueForInvalidToken(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	token, err := utils.GenerateAccessToken(helper.CustomerId.String())
	assert.Nil(t, err)

	resultVerify, err := utils.VerifyAccessToken(token)
	assert.Nil(t, err)

	assert.True(t, resultVerify)
}

func TestVerifyAccessToken_ShouldReturnFalseForInvalidToken(t *testing.T) {
	utils.InitJwtConfig([]byte("abc"), 10)

	resultVerify, err := utils.VerifyAccessToken("Invalid-token")
	assert.NotNil(t, err)
	assert.False(t, resultVerify)
}

func TestExtractIDFromToken_InvalidToken(t *testing.T) {
	tokenString := "invalid token"
	_, err := utils.ExtractIDFromToken(tokenString)
	assert.Error(t, err)
}
