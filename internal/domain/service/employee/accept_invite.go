package employee

import (
	"context"
	"errors"
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
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}

	if err = s.jwt.VerifyInviteToken(token, inv.Token); err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invite token mismatch"),
		)
	}

	_, err = s.Get(ctx, GetFilters{
		UserID: &userID,
	})
	if err == nil {
		return models.Invite{}, errx.ErrorEmployeeAlreadyExists.Raise(
			fmt.Errorf("user is already an employee"),
		)
	}
	if !errors.Is(err, errx.ErrorEmployeeNotFound) {
		return models.Invite{}, err
	}

	err = s.DistributorIsActive(ctx, inv.DistributorID)
	if err != nil {
		return models.Invite{}, err
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		_, err = s.Create(ctx, CreateParams{
			UserID:        userID,
			DistributorID: inv.DistributorID,
			Role:          inv.Role,
		})
		if err != nil {
			return err
		}

		if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusAccepted, now); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("update invite status: %w", err),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = &userID
	inv.AnsweredAt = &now

	return inv, nil
}
