package entity

import (
	"github.com/google/uuid"
)

type Merchant struct {
	Id        uuid.UUID
	Name      string
	CreatedAt string
	UpdatedAt string
}
