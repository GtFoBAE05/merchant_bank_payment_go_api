package main

import (
	"log"
	"merchant_bank_payment_go_api/internal/config"
	"merchant_bank_payment_go_api/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	utils.InitJwtConfig(cfg.SecretKey, cfg.ExpireInMinutes)

	logger := config.NewLogger()

	router := config.Bootstrap(logger)

	if err := router.Run(":4000"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
