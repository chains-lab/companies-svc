package models

import (
	"time"

	"github.com/google/uuid"
)

type Block struct {
	ID            uuid.UUID  `json:"id"`
	DistributorID uuid.UUID  `json:"distributor_id"`
	InitiatorID   uuid.UUID  `json:"initiator_id"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"` // e.g., "active", "canceled"
	BlockedAt     time.Time  `json:"blocked_at"`
	CanceledAt    *time.Time `json:"canceled_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
