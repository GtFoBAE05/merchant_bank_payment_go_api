package repository

type AuthRepository interface {
	LoadBlacklist() ([]string, error)
	SaveBlacklist(blacklistedTokens []string) error
	AddToBlacklist(token string) error
	IsTokenBlacklisted(token string) (bool, error)
}
