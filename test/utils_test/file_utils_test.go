package utils_test

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/utils"
	"merchant_bank_payment_go_api/test/helper"
	"os"
	"testing"
)

var structString = fmt.Sprintf("%q", helper.ExpectedCustomers)

func TestReadJsonFile_Success(t *testing.T) {
	log := logrus.New()
	expectedContent := []byte(structString)

	err := os.WriteFile(helper.FileUtilsFileName, expectedContent, 0644)
	if err != nil {
		t.Errorf("Error creating file %s: %v", helper.FileUtilsFileName, err)
	}
	defer func(filename string) {
		errDefer := os.Remove(filename)
		if errDefer != nil {
			t.Errorf("error on remove file: %v", err)
		}
	}(helper.FileUtilsFileName)

	contentResult, err := utils.ReadJsonFile(helper.FileUtilsFileName, log)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, expectedContent, contentResult)
}

func TestReadJsonFile_FileNotFound(t *testing.T) {
	log := logrus.New()
	contentResult, err := utils.ReadJsonFile(helper.FileUtilsFileName, log)

	assert.Nil(t, contentResult)
	assert.NotNil(t, err)
}

func TestReadJsonFile_ReadError(t *testing.T) {
	log := logrus.New()
	err := os.WriteFile(helper.FileUtilsFileName, []byte(structString), 0000)
	if err != nil {
		t.Errorf("Error creating file %s: %v", helper.FileUtilsFileName, err)
	}
	defer func(filename string) {
		errDefer := os.Remove(filename)
		if errDefer != nil {
			t.Errorf("error on remove file: %v", err)
		}
	}(helper.FileUtilsFileName)

	content, err := utils.ReadJsonFile(helper.FileUtilsFileName, log)

	assert.Nil(t, content)
	assert.NotNil(t, err)
}

func TestWriteJsonFile_Success(t *testing.T) {
	log := logrus.New()
	err := utils.WriteJsonFile(helper.FileUtilsFileName, helper.ExpectedCustomers, log)
	if err != nil {
		t.Errorf("Error write file %s: %v", helper.FileUtilsFileName, err)
	}

	fileContent, err := os.ReadFile(helper.FileUtilsFileName)
	if err != nil {
		t.Errorf("Error read file %s: %v", helper.FileUtilsFileName, err)
	}
	defer func(filename string) {
		errDefer := os.Remove(filename)
		if errDefer != nil {
			t.Errorf("error on remove file: %v", err)
		}
	}(helper.FileUtilsFileName)

	var result []entity.Customer
	err = json.Unmarshal(fileContent, &result)
	if err != nil {
		t.Errorf("failed to unmarshal file content: %v", err)
	}

	assert.Equal(t, len(helper.ExpectedCustomers), len(result))
}

func TestWriteJsonFile_EncodingError(t *testing.T) {
	log := logrus.New()
	testData := func() {}

	err := utils.WriteJsonFile(helper.FileUtilsFileName, testData, log)

	assert.NotNil(t, err)

	defer func(filename string) {
		errDefer := os.Remove(filename)
		if errDefer != nil {
			t.Errorf("error on remove file: %v", err)
		}
	}(helper.FileUtilsFileName)
}
