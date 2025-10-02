package models

import (
	"time"

	"github.com/google/uuid"
)

type DistributorBlock struct {
	ID            uuid.UUID  `json:"id"`
	DistributorID uuid.UUID  `json:"distributor_id"`
	InitiatorID   uuid.UUID  `json:"initiator_id"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"` // e.g., "active", "canceled"
	BlockedAt     time.Time  `json:"blocked_at"`
	CanceledAt    *time.Time `json:"canceled_at"`
}

func (b DistributorBlock) IsNil() bool {
	return b.ID == uuid.Nil
}

type DistributorBlockCollection struct {
	Items []DistributorBlock `json:"items"`
	Page  uint64             `json:"page"`
	Size  uint64             `json:"size"`
	Total uint64             `json:"total"`
}
