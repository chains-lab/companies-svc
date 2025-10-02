package employee

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error) {
	data, err := s.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid or expired token: %w", err),
		)
	}

	now := time.Now().UTC()

	inv, err := s.GetInvite(ctx, data.InviteID)
	if err != nil {
		return models.Invite{}, err
	}
	if inv.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(fmt.Errorf("invite expired"))
	}
	if data.DistributorID != inv.DistributorID {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(fmt.Errorf("token distributor_id mismatch"))
	}

	_, err = s.CreateEmployee(ctx, CreateParams{
		UserID:        userID,
		DistributorID: inv.DistributorID,
		Role:          inv.Role,
	})
	if err != nil {
		return models.Invite{}, err
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = &userID
	inv.AnsweredAt = &now

	if err = s.invite.New().FilterID(inv.ID).Update(ctx, data.UpdateInviteParams{
		Status:     &inv.Status,
		UserID:     &uuid.NullUUID{UUID: userID, Valid: true},
		AnsweredAt: &sql.NullTime{Time: now, Valid: true},
	}); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("update invite status: %w", err),
		)
	}

	return inv, nil
}
