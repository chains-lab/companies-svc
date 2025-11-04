package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Role      string    `json:"role"`
	Position  *string   `json:"position"`
	Label     *string   `json:"label,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e Employee) IsNil() bool {
	return e.UserID == uuid.Nil
}

type EmployeeCollection struct {
	Data  []Employee `json:"data"`
	Page  uint64     `json:"page"`
	Size  uint64     `json:"size"`
	Total uint64     `json:"total"`
}
