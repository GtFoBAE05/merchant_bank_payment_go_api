package usecase

type HistoryUseCase interface {
	AddHistory(customerId, action, details string) error
}
