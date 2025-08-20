package models

import (
	"time"

	"github.com/google/uuid"
)

type Distributor struct {
	ID        uuid.UUID `json:"id"`
	Icon      string    `json:"icon"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
