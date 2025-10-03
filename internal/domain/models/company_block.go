package models

import (
	"time"

	"github.com/google/uuid"
)

type CompanyBlock struct {
	ID          uuid.UUID  `json:"id"`
	CompanyID   uuid.UUID  `json:"company_id"`
	InitiatorID uuid.UUID  `json:"initiator_id"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"` // e.g., "active", "canceled"
	BlockedAt   time.Time  `json:"blocked_at"`
	CanceledAt  *time.Time `json:"canceled_at"`
}

func (b CompanyBlock) IsNil() bool {
	return b.ID == uuid.Nil
}

type CompanyBlockCollection struct {
	Data  []CompanyBlock `json:"data"`
	Page  uint64         `json:"page"`
	Size  uint64         `json:"size"`
	Total uint64         `json:"total"`
}
