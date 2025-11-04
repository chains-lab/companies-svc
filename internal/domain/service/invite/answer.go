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

func (s Service) Answer(ctx context.Context, userID, inviteID uuid.UUID, answer string) (models.Invite, error) {
	err := enum.CheckInviteStatus(answer)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteStatus.Raise(
			fmt.Errorf("invalid invite answer: %w", err),
		)
	}

	now := time.Now().UTC()

	inv, err := s.GetInvite(ctx, inviteID)
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

	if inv.UserID != userID {
		return models.Invite{}, errx.ErrorInviteNotForUser.Raise(
			fmt.Errorf("invite not for user %s", userID),
		)
	}

	emp, err := s.db.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user_id %s, cause: %w", userID, err),
		)
	}
	if !emp.IsNil() {
		return models.Invite{}, errx.ErrorUserAlreadyEmployee.Raise(
			fmt.Errorf("employee with user_id %s already exists", userID),
		)
	}

	err = s.companyIsActive(ctx, inv.CompanyID)
	if err != nil {
		return models.Invite{}, err
	}

	employee := models.Employee{
		UserID:    userID,
		CompanyID: inv.CompanyID,
		Role:      inv.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CreateEmployee(ctx, employee)
		if err != nil {
			return err
		}

		if err = s.db.UpdateInviteStatus(ctx, inv.ID, enum.InviteStatusAccepted); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Invite{}, err
	}

	if err = s.event.PublishEmployeeCreated(ctx, employee); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee with event, cause: %w", err),
		)
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = userID

	return inv, nil
}
