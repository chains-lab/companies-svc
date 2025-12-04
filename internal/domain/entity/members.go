package entity

import (
	"fmt"
	"time"

	"github.com/chains-lab/organizations-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type Member struct {
	UserID    uuid.UUID `json:"user_id"`
	OrgID     uuid.UUID `json:"organization_id"`
	RoleID    uuid.UUID `json:"role_id"`
	Position  *string   `json:"position,omitempty"`
	Label     *string   `json:"label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateCreateMember(mem Member) error {
	switch {
	case mem.UserID == uuid.Nil:
		return errx.ErrMemberInvalidData.Raise(fmt.Errorf("user id should not be empty"))
	case mem.OrgID == uuid.Nil:
		return errx.ErrMemberInvalidData.Raise(fmt.Errorf("organization id should not be empty"))
	case mem.RoleID == uuid.Nil:
		return errx.ErrMemberInvalidData.Raise(fmt.Errorf("role id should not be empty"))
	case mem.CreatedAt == time.Time{}:
		return errx.ErrMemberInvalidData.Raise(fmt.Errorf("created at should not be empty"))
	default:
		return nil
	}
}

func (m Member) IsNil() bool {
	if m.UserID == uuid.Nil {
		return true
	}

	return false
}
