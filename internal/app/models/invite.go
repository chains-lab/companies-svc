package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID            uuid.UUID  `json:"id"`
	DistributorID uuid.UUID  `json:"distributor_id"`
	UserID        uuid.UUID  `json:"user_id"`
	InvitedBy     uuid.UUID  `json:"invited_by"`
	Role          string     `json:"role"` // enum employee_roles
	Status        string     `json:"status"`
	AnsweredAt    *time.Time `json:"answered_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
