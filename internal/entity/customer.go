package entity

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
