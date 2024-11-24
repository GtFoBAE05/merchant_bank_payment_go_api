package entity

import (
	"github.com/google/uuid"
)

type History struct {
	Id         uuid.UUID
	Action     string
	CustomerId uuid.UUID
	Timestamp  string
	Details    string
}
