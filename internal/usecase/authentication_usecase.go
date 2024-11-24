package usecase

import "merchant_bank_payment_go_api/internal/model"

type AuthUseCase interface {
	Login(request model.LoginRequest) (model.LoginResponse, error)
	Logout(accessToken string) error
	IsTokenBlacklisted(accessToken string) (bool, error)
	AddToBlacklist(accessToken string) error
}
