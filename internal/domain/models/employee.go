package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Role      string    `json:"role"`
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

type EmployeeWithUserData struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    *string   `json:"avatar"`
	CompanyID uuid.UUID `json:"company_id"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e EmployeeWithUserData) IsNil() bool {
	return e.UserID == uuid.Nil
}

func (e Employee) AddProfileData(profile Profile) EmployeeWithUserData {
	return EmployeeWithUserData{
		UserID:    e.UserID,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		CompanyID: e.CompanyID,
		Role:      e.Role,
		UpdatedAt: e.UpdatedAt,
		CreatedAt: e.CreatedAt,
	}
}

type EmployeeWithUserDataCollection struct {
	Data  []EmployeeWithUserData `json:"data"`
	Page  uint64                 `json:"page"`
	Size  uint64                 `json:"size"`
	Total uint64                 `json:"total"`
}

func (c EmployeeCollection) AddProfileData(profiles map[uuid.UUID]Profile) EmployeeWithUserDataCollection {
	employees := make([]EmployeeWithUserData, 0, len(c.Data))
	for _, emp := range c.Data {
		empWithProfile := emp.AddProfileData(profiles[emp.UserID])
		employees = append(employees, empWithProfile)
	}
	return EmployeeWithUserDataCollection{
		Data:  employees,
		Page:  c.Page,
		Size:  c.Size,
		Total: c.Total,
	}
}
