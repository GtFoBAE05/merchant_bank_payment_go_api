package usecase_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/usecase/impl"
	"merchant_bank_payment_go_api/test/helper"
	"testing"
)

func TestAddHistory_ShouldCallRepository(t *testing.T) {
	customerId := uuid.New()
	action := "LOGIN"
	details := "Login successful"

	mockHistoryRepository := new(helper.MockHistoryRepository)
	mockHistoryRepository.On("AddHistory", mock.MatchedBy(func(h entity.History) bool {
		return h.CustomerId == customerId &&
			h.Action == action &&
			h.Details == details
	})).Return(nil)

	log := logrus.New()
	historyUseCase := impl.NewHistoryUseCaseImpl(log, mockHistoryRepository)

	err := historyUseCase.AddHistory(customerId.String(), action, details)

	assert.Nil(t, err)
	mockHistoryRepository.AssertExpectations(t)
}

func TestAddHistory_ShouldReturnErrorWhenInvalidCustomerId(t *testing.T) {
	log := logrus.New()

	mockHistoryRepository := new(helper.MockHistoryRepository)
	authUseCase := impl.NewHistoryUseCaseImpl(log, mockHistoryRepository)

	err := authUseCase.AddHistory("invalid_id", "LOGIN", "Login successful")

	assert.NotNil(t, err)
	mockHistoryRepository.AssertExpectations(t)
}

func TestLogAndAddHistory_ShouldReturnNilWhenErrorStatusIsNil(t *testing.T) {
	customerId := uuid.New()
	action := "LOGIN"
	details := "Login successful"

	mockHistoryRepository := new(helper.MockHistoryRepository)
	mockHistoryRepository.On("AddHistory", mock.MatchedBy(func(h entity.History) bool {
		return h.CustomerId == customerId &&
			h.Action == action &&
			h.Details == details
	})).Return(nil)

	log := logrus.New()
	authUseCase := impl.NewHistoryUseCaseImpl(log, mockHistoryRepository)

	err := authUseCase.LogAndAddHistory(customerId.String(), action, details, nil)
	assert.Nil(t, err)
}

func TestLogAndAddHistory_ShouldReturnNilWhenErrorStatusIsNotNil(t *testing.T) {
	customerId := uuid.New()
	action := "LOGIN"
	details := "Login failed"

	mockHistoryRepository := new(helper.MockHistoryRepository)
	mockHistoryRepository.On("AddHistory", mock.MatchedBy(func(h entity.History) bool {
		return h.CustomerId == customerId &&
			h.Action == action &&
			h.Details == details
	})).Return(nil)

	log := logrus.New()
	authUseCase := impl.NewHistoryUseCaseImpl(log, mockHistoryRepository)

	err := authUseCase.LogAndAddHistory(customerId.String(), action, details, errors.New("something wrong"))
	assert.Nil(t, err)
}

func TestLogAndAddHistory_ShouldReturnErrorWhenAddHistoryError(t *testing.T) {
	customerId := uuid.New()
	action := "LOGIN"
	details := "Login failed"

	mockHistoryRepository := new(helper.MockHistoryRepository)
	mockHistoryRepository.On("AddHistory", mock.MatchedBy(func(h entity.History) bool {
		return h.CustomerId == customerId &&
			h.Action == action &&
			h.Details == details
	})).Return(nil)

	log := logrus.New()
	authUseCase := impl.NewHistoryUseCaseImpl(log, mockHistoryRepository)

	err := authUseCase.LogAndAddHistory("12345", action, details, errors.New("something wrong"))
	assert.NotNil(t, err)
}
