package entity

import (
	"github.com/google/uuid"
	"time"
)

type History struct {
	Id         uuid.UUID `json:"id"`
	Action     string    `json:"action"`
	CustomerId string    `json:"customer_id"`
	Timestamp  time.Time `json:"timestamp"`
	Details    string    `json:"details"`
}
