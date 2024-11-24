package impl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type AuthUseCaseImpl struct {
	Log             *logrus.Logger
	AuthRepository  repository.AuthRepository
	CustomerUseCase usecase.CustomerUseCase
	HistoryUseCase  usecase.HistoryUseCase
}

func NewAuthUseCaseImpl(log *logrus.Logger, authRepository repository.AuthRepository, customerUseCase usecase.CustomerUseCase, historyUseCase usecase.HistoryUseCase) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		Log:             log,
		AuthRepository:  authRepository,
		CustomerUseCase: customerUseCase,
		HistoryUseCase:  historyUseCase,
	}
}

func (c *AuthUseCaseImpl) Login(request model.LoginRequest) (model.LoginResponse, error) {
	c.Log.Infof("Attempting login for username: %s", request.Username)

	customer, err := c.CustomerUseCase.FindByUsername(request.Username)
	if err != nil {
		errLog := c.HistoryUseCase.LogAndAddHistory("-", "LOGIN", fmt.Sprintf("Login failed because customer with username %s not exists", request.Username), err)
		if errLog != nil {
			c.Log.Errorf("Failed to log history: %v", errLog)
		}
		return model.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password)); err != nil {
		errLog := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Invalid credential", err)
		if errLog != nil {
			c.Log.Errorf("Failed to log history: %v", errLog)
		}
		return model.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	accessToken, err := jwtutils.GenerateAccessToken(customer.Id.String())
	if err != nil {
		errLog := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Failed to generate access token", err)
		if errLog != nil {
			c.Log.Errorf("Failed to log history: %v", errLog)
		}
		return model.LoginResponse{}, err
	}

	errLog := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Login successful", nil)
	if errLog != nil {
		c.Log.Errorf("Failed to log history: %v", errLog)
		return model.LoginResponse{}, errLog
	}

	return model.LoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (c *AuthUseCaseImpl) Logout(accessToken string) error {
	c.Log.Infof("Attempting logout for access token: %s", accessToken)

	userId, err := jwtutils.ExtractIDFromToken(accessToken)
	if err != nil {
		errLog := c.HistoryUseCase.LogAndAddHistory("-", "LOGOUT", fmt.Sprintf("Logout failed: %v", err), err)
		if errLog != nil {
			c.Log.Errorf("Failed to log history: %v", errLog)
		}
		return err
	}

	errLog := c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Customer ID extracted successfully", nil)
	if errLog != nil {
		c.Log.Errorf("Failed to log history: %v", errLog)
		return errLog
	}

	if err := c.AuthRepository.AddToBlacklist(accessToken); err != nil {
		errLog := c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", fmt.Sprintf("Failed to blacklist token: %v", err), err)
		if errLog != nil {
			c.Log.Errorf("Failed to log history: %v", errLog)
		}
		return err
	}

	errLog = c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Token blacklisted successfully", nil)
	if errLog != nil {
		c.Log.Errorf("Failed to log history: %v", errLog)
		return errLog
	}

	errLog = c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Logout successful", nil)
	if errLog != nil {
		c.Log.Errorf("Failed to log history: %v", errLog)
		return errLog
	}

	return nil
}

func (c *AuthUseCaseImpl) IsTokenBlacklisted(accessToken string) (bool, error) {
	return c.AuthRepository.IsTokenBlacklisted(accessToken)
}

func (c *AuthUseCaseImpl) AddToBlacklist(accessToken string) error {
	return c.AuthRepository.AddToBlacklist(accessToken)
}
