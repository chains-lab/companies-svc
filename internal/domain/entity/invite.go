package entity

import (
	"fmt"
	"time"

	"github.com/chains-lab/organizations-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type Invite struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"organization_id"`
	UserID      uuid.UUID `json:"user_id"`
	InitiatorID uuid.UUID `json:"initiator_id"`
	RoleID      uuid.UUID `json:"role_id"`
	Status      string    `json:"status"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func ValidateCreateInvite(inv Invite) error {
	switch {
	case inv.ID == uuid.Nil:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with nil ID"))
	case inv.OrgID == uuid.Nil:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with nil OrgID"))
	case inv.UserID == uuid.Nil:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with nil UserID"))
	case inv.InitiatorID == uuid.Nil:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with nil InitiatorID"))
	case inv.RoleID == uuid.Nil:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with nil RoleID"))
	case ValidateInviteStatus(inv.Status) != nil:
		return ValidateInviteStatus(inv.Status)
	case inv.ExpiresAt.After(time.Now()):
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with incorrect ExpiresAt"))
	case inv.CreatedAt != time.Time{}:
		return errx.ErrorInviteInvalidData.Raise(fmt.Errorf("cannot create invite with incorrect CreatedAt"))
	default:
		return nil
	}
}

func (i Invite) IsNil() bool {
	if i.ID == uuid.Nil {
		return true
	}

	return false
}

func (i Invite) CanBeAnswered() bool {
	if i.Status != InviteStatusSent {
		return false
	}

	return time.Now().UTC().Before(i.ExpiresAt)
}

const (
	InviteStatusSent     = "sent"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
)

var inviteStatuses = []string{
	InviteStatusSent,
	InviteStatusAccepted,
	InviteStatusDeclined,
}

func GetAllInviteStatuses() []string {
	return inviteStatuses
}

var ErrorInviteStatusNotSupported = fmt.Errorf("invite status not supported must be one of: %v", GetAllInviteStatuses())

func ValidateInviteStatus(status string) error {
	for _, s := range inviteStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorInviteStatusNotSupported)
}
