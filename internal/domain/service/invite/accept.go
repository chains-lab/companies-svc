package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Accept(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error) {
	data, err := s.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("failed invalid or expired token, cause: %w", err),
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

	emp, err := s.db.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user_id %s, cause: %w", userID, err),
		)
	}
	if !emp.IsNil() {
		return models.Invite{}, errx.ErrorEmployeeAlreadyExists.Raise(
			fmt.Errorf("employee with user_id %s already exists", userID),
		)
	}

	err = s.companyIsActive(ctx, inv.CompanyID)
	if err != nil {
		return models.Invite{}, err
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CreateEmployee(ctx, models.Employee{
			UserID:    userID,
			CompanyID: inv.CompanyID,
			Role:      inv.Role,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}

		if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusAccepted, now); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Invite{}, err
	}

	if err = s.eve.UpdateEmployee(ctx, userID, &inv.CompanyID, &inv.Role); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee with event, cause: %w", err),
		)
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = &userID
	inv.AnsweredAt = &now

	return inv, nil
}
