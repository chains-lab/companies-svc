package employee

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type AnsweredInviteParams struct {
	Token  string
	Status string
}

func (e Employee) AnsweredInvite(ctx context.Context, userID uuid.UUID, params AnsweredInviteParams) (models.Invite, error) {
	data, err := e.jwt.DecryptInviteToken(params.Token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid or expired token: %w", err),
		)
	}

	inv, err := e.GetInvite(ctx, data.JTI)
	if err != nil {
		return models.Invite{}, err
	}

	now := time.Now().UTC()

	if inv.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(fmt.Errorf("invite expired"))
	}
	if data.DistributorID != inv.DistributorID {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(fmt.Errorf("token city_id mismatch"))
	}

	switch params.Status {
	case enum.InviteStatusAccepted:
		_, err = e.CreateEmployee(ctx, CreateParams{
			UserID:        userID,
			DistributorID: inv.DistributorID,
			Role:          inv.Role,
		})
		if err != nil {
			return models.Invite{}, err
		}
	case enum.InviteStatusRejected:
		// nothing to do
	default:
		return models.Invite{}, errx.ErrorInvalidInviteStatus.Raise(
			fmt.Errorf("invalid invite status: %s", params.Status),
		)
	}

	upd := dbx.UpdateInviteParams{
		Status:     &params.Status,
		UserID:     &uuid.NullUUID{UUID: userID, Valid: true},
		AnsweredAt: &sql.NullTime{Time: now, Valid: true},
	}

	if err = e.invite.New().FilterID(inv.ID).Update(ctx, upd); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("update invite status: %w", err),
		)
	}

	inv.Status = params.Status
	inv.UserID = &userID
	inv.AnsweredAt = &now
	return inv, nil
}
