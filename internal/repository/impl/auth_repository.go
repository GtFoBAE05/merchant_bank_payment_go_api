package impl

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/utils"
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
	r.Log.Debugf("Loading blacklisted tokens from file: %s", r.Filename)

	file, err := utils.ReadJsonFile(r.Filename, r.Log)
	if err != nil {
		r.Log.Errorf("Error reading file %s: %v", r.Filename, err)
		return nil, err
	}

	r.Log.Tracef("File content: %s", string(file))

	var blacklistedTokens []string
	err = json.Unmarshal(file, &blacklistedTokens)
	if err != nil {
		r.Log.Errorf("Error decoding JSON from file %s: %v", r.Filename, err)
		return nil, err
	}

	r.Log.Infof("Successfully loaded %d blacklisted tokens from %s", len(blacklistedTokens), r.Filename)
	return blacklistedTokens, nil
}

func (r *AuthRepositoryImpl) SaveBlacklist(blacklistedTokens []string) error {
	r.Log.Infof("Saving %d blacklisted tokens to file: %s", len(blacklistedTokens), r.Filename)

	err := utils.WriteJSONFile(r.Filename, blacklistedTokens, r.Log)
	if err != nil {
		r.Log.Errorf("Error saving blacklist to file %s: %v", r.Filename, err)
		return fmt.Errorf("error saving blacklist to file %s: %v", r.Filename, err)
	}

	r.Log.Infof("Successfully saved blacklist to %s", r.Filename)
	return nil
}

func (r *AuthRepositoryImpl) AddToBlacklist(token string) error {
	blacklistedTokens, err := r.LoadBlacklist()
	if err != nil {
		return err
	}

	for _, blacklistedToken := range blacklistedTokens {
		if blacklistedToken == token {
			r.Log.Warnf("Token %s is already blacklisted", token)
			return fmt.Errorf("token is already blacklisted")
		}
	}

	r.Log.Infof("Adding token %s to blacklist", token)

	blacklistedTokens = append(blacklistedTokens, token)

	return r.SaveBlacklist(blacklistedTokens)
}

func (r *AuthRepositoryImpl) IsTokenBlacklisted(token string) (bool, error) {
	blacklistedTokens, err := r.LoadBlacklist()
	if err != nil {
		return false, err
	}

	for _, blacklistedToken := range blacklistedTokens {
		if blacklistedToken == token {
			r.Log.Infof("Token %s is blacklisted", token)
			return true, nil
		}
	}

	r.Log.Debugf("Token %s is not blacklisted", token)
	return false, nil
}
