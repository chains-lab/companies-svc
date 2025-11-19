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

func (s Service) Reply(
	ctx context.Context,
	userID, inviteID uuid.UUID,
	reply string,
) (models.Invite, error) {
	err := enum.CheckInviteStatus(reply)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteStatus.Raise(
			fmt.Errorf("invalid invite reply: %w", err),
		)
	}

	now := time.Now().UTC()

	invite, err := s.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyReplyed.Raise(
			fmt.Errorf("invite already replyed with status=%s", invite.Status),
		)
	}
	if now.After(invite.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}
	if invite.UserID != userID {
		return models.Invite{}, errx.ErrorInviteNotForThisUser.Raise(
			fmt.Errorf("invite not for user %s", userID),
		)
	}

	emp, err := s.db.GetEmployee(ctx, userID, invite.CompanyID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf(
				"failed to get employee with user_id %s in company %s, cause: %w",
				userID, invite.CompanyID, err,
			),
		)
	}
	if !emp.IsNil() {
		return models.Invite{}, errx.ErrorUserAlreadyInThisCompany.Raise(
			fmt.Errorf("employee with user_id %s already in company %s", userID, invite.CompanyID),
		)
	}

	company, err := s.getCompany(ctx, invite.CompanyID)
	if err != nil {
		return models.Invite{}, err
	}
	if company.Status != enum.CompanyStatusActive {
		return models.Invite{}, errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("company with EmployeeID %s is not active", invite.CompanyID),
		)
	}

	employee := models.Employee{
		UserID:    userID,
		CompanyID: invite.CompanyID,
		Role:      invite.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	switch reply {
	case enum.InviteStatusAccepted:
		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			if err = s.db.UpdateInviteStatus(ctx, invite.ID, enum.InviteStatusAccepted); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update invite status, cause: %w", err),
				)
			}

			err = s.db.CreateEmployee(ctx, employee)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to create employee, cause: %w", err),
				)
			}

			return nil
		}); err != nil {
			return models.Invite{}, err
		}

		err = s.event.PublishInviteAccepted(ctx, invite, models.Company{ID: invite.CompanyID}, invite.InitiatorID)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite accepted event, cause: %w", err),
			)
		}

		err = s.event.PublishEmployeeCreated(ctx, company, employee, invite.InitiatorID)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish employee created event, cause: %w", err),
			)
		}

	case enum.InviteStatusDeclined:
		if err = s.db.UpdateInviteStatus(ctx, invite.ID, enum.InviteStatusDeclined); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}

		err = s.event.PublishInviteDeclined(ctx, invite, company)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite declined event, cause: %w", err),
			)
		}
	}

	invite.Status = reply

	return invite, nil
}
