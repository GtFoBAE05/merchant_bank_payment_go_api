package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func ReadJsonFile(filename string, log *logrus.Logger) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("Error reading file %s: %v", filename, err)
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}
	return fileContent, nil
}

func WriteJSONFile(filename string, data interface{}, log *logrus.Logger) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating file %s: %v", filename, err)
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Error closing file %s: %v", filename, err)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		log.Errorf("Error encoding data to file %s: %v", filename, err)
		return fmt.Errorf("error encoding data to file %s: %w", filename, err)
	}

	return nil
}
