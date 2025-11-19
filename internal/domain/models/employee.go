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

type EmployeesCollection struct {
	Data  []Employee `json:"data"`
	Page  uint64     `json:"page"`
	Size  uint64     `json:"size"`
	Total uint64     `json:"total"`
}

func (c EmployeesCollection) GetUserIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(c.Data))
	for _, admin := range c.Data {
		ids = append(ids, admin.UserID)
	}
	return ids
}
