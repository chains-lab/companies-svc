package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	DistributorID uuid.UUID
	Role          string
}

func (s Service) CreateInvite(ctx context.Context, InitiatorID uuid.UUID, params SentInviteParams) (models.Invite, error) {
	initiator, err := s.GetInitiator(ctx, InitiatorID)
	if err != nil {
		return models.Invite{}, err
	}

	if initiator.DistributorID != params.DistributorID {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployeeOfThisDistributor.Raise(
			fmt.Errorf("initiator distributor_id %s not equal to params distributor_id %s", initiator.DistributorID, params.DistributorID),
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

	if err = s.DistributorIsActive(ctx, initiator.DistributorID); err != nil {
		return models.Invite{}, err
	}

	inviteID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := s.jwt.CreateInviteToken(inviteID, params.Role, params.DistributorID, exAt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	hash, err := s.jwt.HashInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("hash invite token: %w", err),
		)
	}

	invite := models.Invite{
		ID:            inviteID,
		Status:        enum.InviteStatusSent,
		Role:          params.Role,
		DistributorID: initiator.DistributorID,
		Token:         hash,
		ExpiresAt:     exAt,
		CreatedAt:     now,
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return invite, nil
}
