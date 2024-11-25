package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.InfoLevel)

	// Set the JSONFormatter for structured logging (machine-readable)
	log.SetFormatter(&logrus.JSONFormatter{})

	logFile, err := os.OpenFile("app_history.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file for development, logging to stderr instead:", err)
	} else {
		log.SetOutput(logFile)
	}

	return log
}
