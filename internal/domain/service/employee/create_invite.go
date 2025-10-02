package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/infra/jwtmanager"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	InitiatorID   uuid.UUID
	DistributorID uuid.UUID
	Role          string
}

func (s Service) CreateInvite(ctx context.Context, params SentInviteParams) (models.Invite, error) {
	initiator, err := s.GetInitiator(ctx, params.InitiatorID)
	if err != nil {
		return models.Invite{}, err
	}

	if initiator.DistributorID != params.DistributorID {
		return models.Invite{}, errx.ErrorInitiatorIsNotThisDistributorEmployee.Raise(
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
		return models.Invite{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	invID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := s.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		ID:            invID,
		DistributorID: initiator.DistributorID,
		Role:          params.Role,
		ExpiredAt:     exAt,
		CreatedAt:     now,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	stmt := pgdb.Invite{
		ID:            invID,
		Status:        enum.InviteStatusSent,
		Role:          params.Role,
		DistributorID: initiator.DistributorID,
		ExpiresAt:     exAt,
		CreatedAt:     now,
	}

	err = s.invite.New().Insert(ctx, stmt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return inviteFromDB(stmt, token), nil
}
