package models

import (
	"time"

	"github.com/google/uuid"
)

type Distributor struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (d Distributor) IsNil() bool {
	return d.ID == uuid.Nil
}

type DistributorCollection struct {
	Items []Distributor `json:"items"`
	Page  uint64        `json:"page"`
	Size  uint64        `json:"size"`
	Total uint64        `json:"total"`
}
