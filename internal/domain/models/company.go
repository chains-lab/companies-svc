package models

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (d Company) IsNil() bool {
	return d.ID == uuid.Nil
}

type CompanyCollection struct {
	Data  []Company `json:"data"`
	Page  uint64    `json:"page"`
	Size  uint64    `json:"size"`
	Total uint64    `json:"total"`
}
