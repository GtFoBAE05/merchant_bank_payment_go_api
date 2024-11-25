package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	SecretKey       []byte
	ExpireInMinutes int
	Port            string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("SECRET_KEY not set in environment variables")
	}

	expireInMinutes, err := strconv.Atoi(os.Getenv("EXPIRE_IN_MINUTES"))
	if err != nil || expireInMinutes <= 0 {
		return nil, fmt.Errorf("EXPIRE_IN_MINUTES is not set or invalid")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	return &Config{
		SecretKey:       []byte(secretKey),
		ExpireInMinutes: expireInMinutes,
		Port:            port,
	}, nil
}
