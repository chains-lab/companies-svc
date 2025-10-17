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
	Role      string
}

func (s Service) Create(ctx context.Context, InitiatorID uuid.UUID, params CreateParams) (models.Invite, error) {
	initiator, err := s.db.GetEmployeeByUserID(ctx, InitiatorID)
	if err != nil {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("failed to get initiator employee by user id %s, cause: %w", InitiatorID, err),
		)
	}

	if initiator.IsNil() {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("initiator employee with user id %s not found", InitiatorID),
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
			fmt.Errorf("compare city gov roles: %w", err),
		)
	}
	if access <= 0 {
		return models.Invite{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	if err = s.companyIsActive(ctx, initiator.CompanyID); err != nil {
		return models.Invite{}, err
	}

	inviteID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := s.jwt.CreateInviteToken(inviteID, params.Role, params.CompanyID, exAt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite token, cause: %w", err),
		)
	}

	hash, err := s.jwt.HashInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to hash invite token, cause: %w", err),
		)
	}

	invite := models.Invite{
		ID:        inviteID,
		Status:    enum.InviteStatusSent,
		Role:      params.Role,
		CompanyID: initiator.CompanyID,
		Token:     hash,
		ExpiresAt: exAt,
		CreatedAt: now,
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite in db, cause: %w", err),
		)
	}

	invite.Token = token

	return invite, nil
}
