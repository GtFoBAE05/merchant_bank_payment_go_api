package impl

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/utils"
)

type CustomerRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewCustomerRepositoryImpl(log *logrus.Logger, filename string) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (r *CustomerRepositoryImpl) LoadCustomers() ([]entity.Customer, error) {
	r.Log.Debugf("Loading customers from file: %s", r.Filename)

	file, err := utils.ReadJsonFile(r.Filename, r.Log)
	if err != nil {
		r.Log.Errorf("Error reading file %s: %v", r.Filename, err)
		return nil, err
	}

	var customers []entity.Customer
	err = json.Unmarshal(file, &customers)
	if err != nil {
		r.Log.Errorf("Error decoding JSON from file %s: %v", r.Filename, err)
		return nil, err
	}

	r.Log.Infof("Successfully loaded %d customers from %s", len(customers), r.Filename)
	return customers, nil
}

func (r *CustomerRepositoryImpl) FindById(id uuid.UUID) (entity.Customer, error) {
	r.Log.Debugf("Finding customer by id: %s", id.String())

	customers, err := r.LoadCustomers()
	if err != nil {
		r.Log.Errorf("Error loading customers from file %s: %v", r.Filename, err)
		return entity.Customer{}, err
	}

	for _, customer := range customers {
		if customer.Id == id {
			r.Log.Infof("Found customer with id: %s", id.String())
			return customer, nil
		}
	}

	err = fmt.Errorf("customer with id %s not found in %s", id, r.Filename)
	r.Log.Errorf(err.Error())
	return entity.Customer{}, err
}

func (r *CustomerRepositoryImpl) FindByUsername(username string) (entity.Customer, error) {
	r.Log.Debugf("Finding customer by username: %s", username)

	customers, err := r.LoadCustomers()
	if err != nil {
		r.Log.Errorf("Error loading customers from file %s: %v", r.Filename, err)
		return entity.Customer{}, err
	}

	for _, customer := range customers {
		if customer.Username == username {
			r.Log.Infof("Found customer with username: %s", username)
			return customer, nil
		}
	}

	err = fmt.Errorf("customer with username %s not found in %s", username, r.Filename)
	r.Log.Errorf(err.Error())
	return entity.Customer{}, err
}
