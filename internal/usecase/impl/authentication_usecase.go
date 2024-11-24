package impl

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type AuthUseCaseImpl struct {
	AuthRepository  repository.AuthRepository
	CustomerUseCase usecase.CustomerUseCase
	HistoryUseCase  usecase.HistoryUseCase
}

func NewAuthUseCaseImpl(authRepository repository.AuthRepository, customerUseCase usecase.CustomerUseCase, historyUseCase usecase.HistoryUseCase) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		AuthRepository:  authRepository,
		CustomerUseCase: customerUseCase,
		HistoryUseCase:  historyUseCase,
	}
}

func (c *AuthUseCaseImpl) Login(request model.LoginRequest) (model.LoginResponse, error) {
	customer, err := c.CustomerUseCase.FindByUsername(request.Username)
	if err != nil {
		errLogHistory := c.HistoryUseCase.LogAndAddHistory("-", "LOGIN", fmt.Sprintf("Login failed because customer with username %s not exists", request.Username), err)
		if errLogHistory != nil {
			return model.LoginResponse{}, errLogHistory
		}
		return model.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password))
	if err != nil {
		errLogHistory := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Invalid credentials", err)
		if errLogHistory != nil {
			return model.LoginResponse{}, errLogHistory
		}
		return model.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	accessToken, err := jwtutils.GenerateAccessToken(customer.Id.String())
	if err != nil {
		errLogHistory := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Failed to generate access token", err)
		if errLogHistory != nil {
			return model.LoginResponse{}, errLogHistory
		}
		return model.LoginResponse{}, err
	}

	errLogHistory := c.HistoryUseCase.LogAndAddHistory(customer.Id.String(), "LOGIN", "Login successful", nil)
	if errLogHistory != nil {
		return model.LoginResponse{}, errLogHistory
	}

	return model.LoginResponse{AccessToken: accessToken}, nil
}

func (c *AuthUseCaseImpl) Logout(accessToken string) error {
	userId, err := jwtutils.ExtractIDFromToken(accessToken)
	if err != nil {
		errLogHistory := c.HistoryUseCase.LogAndAddHistory("-", "LOGOUT", fmt.Sprintf("Logout failed: %v", err), err)
		if errLogHistory != nil {
			return errLogHistory
		}
		return err
	}

	errLogHistory := c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Customer ID extracted successfully", nil)
	if errLogHistory != nil {
		return errLogHistory
	}

	err = c.AuthRepository.AddToBlacklist(accessToken)
	if err != nil {
		errLogHistory = c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", fmt.Sprintf("Failed to blacklist token: %v", err), err)
		if errLogHistory != nil {
			return errLogHistory
		}
		return err
	}

	errLogHistory = c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Token blacklisted successfully", nil)
	if errLogHistory != nil {
		return errLogHistory
	}

	errLogHistory = c.HistoryUseCase.LogAndAddHistory(userId, "LOGOUT", "Logout successful", nil)
	if errLogHistory != nil {
		return errLogHistory
	}

	return nil
}

func (c *AuthUseCaseImpl) IsTokenBlacklisted(accessToken string) (bool, error) {
	return c.AuthRepository.IsTokenBlacklisted(accessToken)
}

func (c *AuthUseCaseImpl) AddToBlacklist(accessToken string) error {
	return c.AuthRepository.AddToBlacklist(accessToken)
}
