package entity

import (
	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	MerchantId uuid.UUID
	Amount     int64
	Timestamp  string
}
