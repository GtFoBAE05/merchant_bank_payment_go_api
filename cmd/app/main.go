package main

import (
	"merchant_bank_payment_go_api/internal/config"
)

func main() {
	log := config.NewLogger()

	router := config.Bootstrap(&config.BootstrapConfig{
		Log: log,
	})

	if err := router.Run(":4000"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
