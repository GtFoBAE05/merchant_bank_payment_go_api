package repository

import "merchant_bank_payment_go_api/internal/entity"

type HistoryRepository interface {
	LoadHistories() ([]entity.History, error)
	SaveHistories(histories []entity.History) error
	AddHistory(history entity.History) error
}
