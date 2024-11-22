package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func ReadJsonFile(filename string, log *logrus.Logger) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Errorf("Error opening file %s: %v", filename, err)
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Error closing file %s: %v", filename, err)
		}
	}()

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("Error reading file %s: %v", filename, err)
		return nil, err
	}
	return fileContent, nil
}

func WriteJSONFile(filename string, data interface{}, log *logrus.Logger) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating file %s: %v", filename, err)
		return fmt.Errorf("error creating file %s: %v", filename, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Errorf("Error closing file %s: %v", filename, err)
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		log.Errorf("Error encoding data to file %s: %v", filename, err)
		return fmt.Errorf("error encoding data to file %s: %v", filename, err)
	}

	return nil
}
