package impl

import (
	"github.com/sirupsen/logrus"
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
	return model.LoginResponse{}, nil
}

func (c *AuthUseCase) Logout(accessToken string) error {
	return nil
}

func (c *AuthUseCase) IsTokenBlacklisted(accessToken string) (bool, error) {
	return false, nil
}

func (c *AuthUseCase) AddToBlacklist(accessToken string) error {
	return c.AuthRepository.AddToBlacklist(accessToken)
}
