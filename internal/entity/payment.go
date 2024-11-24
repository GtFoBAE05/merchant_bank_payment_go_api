package entity

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	Id         uuid.UUID `json:"id"`
	CustomerId uuid.UUID `json:"customer_id"`
	MerchantId uuid.UUID `json:"merchant_id"`
	Amount     int64     `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
