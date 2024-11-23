package entity

import "github.com/google/uuid"

type Customer struct {
	Id        uuid.UUID
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}
