package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type JwtConfig struct {
	SecretKey       []byte
	ExpireInMinutes int
}

var jwtConfig *JwtConfig

func InitJwtConfig(secretKey []byte, expireInMinutes int) {
	if len(secretKey) == 0 {
		secretKey = []byte("supersecretkey")
		log.Warn("Jwt Secret Key is missing. Using default value: 'supersecretkey'")
	}

	if expireInMinutes <= 0 {
		expireInMinutes = 10
		log.Warn("Jwt ExpireInMinutes is invalid or missing. Using default value: 10 minutes")
	}

	jwtConfig = &JwtConfig{
		SecretKey:       secretKey,
		ExpireInMinutes: expireInMinutes,
	}
}

func GenerateAccessToken(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(jwtConfig.ExpireInMinutes) * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = id

	tokenString, err := token.SignedString(jwtConfig.SecretKey)
	if err != nil {
		return "", fmt.Errorf("error signing the token: %w", err)
	}
	return tokenString, nil
}

func VerifyAccessToken(accessToken string) (bool, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtConfig.SecretKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	return true, nil
}

func ExtractIDFromToken(requestToken string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtConfig.SecretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid or malformed token")
	}

	id, ok := claims["user"].(string)
	if !ok {
		return "", errors.New("user ID missing or invalid in token")
	}

	return id, nil
}
