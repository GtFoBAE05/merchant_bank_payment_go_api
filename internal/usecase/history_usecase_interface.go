package usecase

type HistoryUseCase interface {
	AddHistory(customerId, action, details string) error
	LogAndAddHistory(userId, action, message string, err error) error
}
