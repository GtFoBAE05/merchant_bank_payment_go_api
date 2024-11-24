package impl

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"merchant_bank_payment_go_api/internal/entity"
	"merchant_bank_payment_go_api/internal/utils"
)

type HistoryRepositoryImpl struct {
	Log      *logrus.Logger
	Filename string
}

func NewHistoryRepositoryImpl(log *logrus.Logger, filename string) *HistoryRepositoryImpl {
	return &HistoryRepositoryImpl{
		Log:      log,
		Filename: filename,
	}
}

func (h *HistoryRepositoryImpl) LoadHistories() ([]entity.History, error) {
	h.Log.Debugf("Loading histories from file: %s", h.Filename)

	file, err := utils.ReadJsonFile(h.Filename, h.Log)
	if err != nil {
		h.Log.Errorf("Error reading file %s: %v", h.Filename, err)
		return nil, err
	}

	h.Log.Tracef("File content: %s", string(file))

	var histories []entity.History
	err = json.Unmarshal(file, &histories)
	if err != nil {
		h.Log.Errorf("Error decoding JSON from file %s: %v", h.Filename, err)
		return nil, err
	}

	h.Log.Infof("Successfully loaded %d histories from %s", len(histories), h.Filename)
	return histories, nil
}

func (h *HistoryRepositoryImpl) SaveHistories(histories []entity.History) error {
	h.Log.Infof("Saving %d histories to file: %s", len(histories), h.Filename)

	err := utils.WriteJSONFile(h.Filename, histories, h.Log)
	if err != nil {
		h.Log.Errorf("Error saving histories to file %s: %v", h.Filename, err)
		return fmt.Errorf("error saving histories to file %s: %v", h.Filename, err)
	}

	h.Log.Infof("Successfully saved histories to %s", h.Filename)
	return nil
}

func (h *HistoryRepositoryImpl) AddHistory(history entity.History) error {
	histories, err := h.LoadHistories()
	if err != nil {
		return err
	}

	h.Log.Infof("Adding history %s to histories", history)

	histories = append(histories, history)

	return h.SaveHistories(histories)
}
