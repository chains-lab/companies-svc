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

type CreateParams struct {
	CompanyID uuid.UUID
	UserID    uuid.UUID
	Role      string
}

func (s Service) Create(ctx context.Context, initiatorID uuid.UUID, params CreateParams) (models.Invite, error) {
	initiator, err := s.db.GetEmployeeByUserID(ctx, initiatorID)
	if err != nil {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("failed to get initiator employee by user id %s, cause: %w", initiatorID, err),
		)
	}
	if initiator.IsNil() {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("initiator employee with user id %s not found", initiatorID),
		)
	}
	if initiator.CompanyID != params.CompanyID {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployeeOfThisCompany.Raise(
			fmt.Errorf("initiator company_id %s not equal to params company_id %s", initiator.CompanyID, params.CompanyID),
		)
	}

	access, err := enum.CompareEmployeeRoles(initiator.Role, params.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("compare employee roles: %w", err),
		)
	}
	if access <= 0 {
		return models.Invite{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	exist, err := s.db.EmployeeExist(ctx, params.UserID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check user existence by user id %s, cause: %w", params.UserID, err),
		)
	}
	if exist {
		return models.Invite{}, errx.ErrorUserAlreadyEmployee.Raise(
			fmt.Errorf("user with id %s not found", params.UserID),
		)
	}

	if err = s.companyIsActive(ctx, initiator.CompanyID); err != nil {
		return models.Invite{}, err
	}

	inviteID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	invite := models.Invite{
		ID:        inviteID,
		CompanyID: initiator.CompanyID,
		UserID:    params.UserID,
		Status:    enum.InviteStatusSent,
		Role:      params.Role,
		ExpiresAt: exAt,
		CreatedAt: now,
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite in db, cause: %w", err),
		)
	}

	err = s.event.PublishInviteCreated(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish invite created event, cause: %w", err),
		)
	}

	return invite, nil
}
