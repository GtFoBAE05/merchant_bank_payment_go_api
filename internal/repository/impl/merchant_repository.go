package impl

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/utils"
)

type MerchantRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewMerchantRepositoryImpl(log *logrus.Logger, filename string) *MerchantRepositoryImpl {
	return &MerchantRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (m *MerchantRepositoryImpl) LoadMerchants() ([]entity.Merchant, error) {
	m.Log.Debugf("Loading merchants from file: %s", m.Filename)

	file, err := utils.ReadJsonFile(m.Filename, m.Log)
	if err != nil {
		m.Log.Errorf("Error reading file %s: %v", m.Filename, err)
		return nil, err
	}

	var merchants []entity.Merchant
	err = json.Unmarshal(file, &merchants)
	if err != nil {
		m.Log.Errorf("Error decoding JSON from file %s: %v", m.Filename, err)
		return nil, err
	}

	m.Log.Infof("Successfully loaded %d merchants from %s", len(merchants), m.Filename)
	return merchants, nil
}

func (m *MerchantRepositoryImpl) FindById(id uuid.UUID) (entity.Merchant, error) {
	m.Log.Debugf("Finding merchant by id: %s", id.String())

	merchants, err := m.LoadMerchants()
	if err != nil {
		m.Log.Errorf("Error loading merchants from file %s: %v", m.Filename, err)
		return entity.Merchant{}, err
	}

	for _, merchant := range merchants {
		if merchant.Id == id {
			m.Log.Infof("Found merchant with id: %s", id.String())
			return merchant, nil
		}
	}

	err = fmt.Errorf("merchant with id %s not found in %s", id, m.Filename)
	m.Log.Errorf(err.Error())
	return entity.Merchant{}, err
}
