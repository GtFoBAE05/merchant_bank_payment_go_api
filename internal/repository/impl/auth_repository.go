package impl

import (
	"github.com/sirupsen/logrus"
)

type AuthRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewAuthRepository(log *logrus.Logger, filename string) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (r *AuthRepositoryImpl) LoadBlacklist() ([]string, error) {
	return nil, nil
}

func (r *AuthRepositoryImpl) SaveBlacklist(blacklistedTokens []string) error {
	return nil
}

func (r *AuthRepositoryImpl) AddToBlacklist(token string) error {
	return r.SaveBlacklist(nil)
}

func (r *AuthRepositoryImpl) IsTokenBlacklisted(token string) (bool, error) {
	return false, nil
}
