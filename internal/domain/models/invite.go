package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	UserID    uuid.UUID
	Role      string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
