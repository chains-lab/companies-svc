package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	UserID        uuid.UUID `json:"user_id"`
	DistributorID uuid.UUID `json:"distributor_id"`
	Role          string    `json:"role"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}
