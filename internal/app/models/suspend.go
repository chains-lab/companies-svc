package models

import (
	"time"

	"github.com/google/uuid"
)

type SuspendedDistributor struct {
	ID            uuid.UUID  `json:"id"`
	DistributorID uuid.UUID  `json:"distributor_id"`
	InitiatorID   uuid.UUID  `json:"initiator_id"`
	Reason        string     `json:"reason"`
	Active        bool       `json:"active"`
	SuspendedAt   time.Time  `json:"suspended_at"`
	CanceledAt    *time.Time `json:"canceled_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
