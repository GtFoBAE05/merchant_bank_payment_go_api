package impl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	auth "merchant_bank_payment_go_api/internal/jwt"
	"merchant_bank_payment_go_api/internal/model"
	"merchant_bank_payment_go_api/internal/repository"
	"merchant_bank_payment_go_api/internal/usecase"
)

type AuthUseCase struct {
	Log             *logrus.Logger
	AuthRepository  repository.AuthRepository
	CustomerUseCase usecase.CustomerUseCaseInterface
}

func NewAuthUseCase(log *logrus.Logger, authRepository repository.AuthRepository, customerUseCase usecase.CustomerUseCaseInterface) *AuthUseCase {
	return &AuthUseCase{
		Log:             log,
		AuthRepository:  authRepository,
		CustomerUseCase: customerUseCase,
	}
}

func (c *AuthUseCase) Login(request model.LoginRequest) (model.LoginResponse, error) {
	c.Log.Infof("Attempting login for username: %s", request.Username)
	customer, err := c.CustomerUseCase.FindByUsername(request.Username)
	if err != nil {
		c.Log.Errorf("Failed to find customer with username %s: %v", request.Username, err)
		return model.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password))
	if err != nil {
		c.Log.Warnf("Invalid password for username %s", request.Username)
		return model.LoginResponse{}, fmt.Errorf("invalid username or password")
	}

	accessToken, err := auth.GenerateAccessToken(customer.Id.String())
	if err != nil {
		c.Log.Errorf("Failed to generate accessToken for user %s: %v", request.Username, err)
		return model.LoginResponse{}, err
	}

	c.Log.Infof("Login successful for username: %s", request.Username)
	return model.LoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (c *AuthUseCase) Logout(accessToken string) error {
	c.Log.Infof("Attempting logout for accessToken: %s", accessToken)

	err := c.AuthRepository.AddToBlacklist(accessToken)
	if err != nil {
		c.Log.Errorf("Failed to blacklist token %s: %v", accessToken, err)
		return err
	}

	c.Log.Infof("Successfully blacklisted token %s", accessToken)
	return nil
}

func (c *AuthUseCase) IsTokenBlacklisted(accessToken string) (bool, error) {
	isBlacklisted, err := c.AuthRepository.IsTokenBlacklisted(accessToken)
	if err != nil {
		c.Log.Errorf("Error checking if token %s is blacklisted: %v", accessToken, err)
	}
	return isBlacklisted, err
}

func (c *AuthUseCase) AddToBlacklist(accessToken string) error {
	c.Log.Infof("Attempting to add token %s to the blacklist", accessToken)
	return c.AuthRepository.AddToBlacklist(accessToken)
}
