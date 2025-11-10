package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	UserID      uuid.UUID `json:"user_id"`
	InitiatorID uuid.UUID `json:"initiator_id"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}
