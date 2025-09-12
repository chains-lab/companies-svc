package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/jwtmanager"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	InitiatorID uuid.UUID
	Role        string
}

func (e Employee) SentInvite(ctx context.Context, params SentInviteParams) (models.Invite, error) {
	initiator, err := e.GetInitiator(ctx, params.InitiatorID)
	if err != nil {
		return models.Invite{}, err
	}

	access, err := enum.CompareCityGovRoles(params.Role, initiator.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("compare city gov roles: %w", err),
		)
	}
	if access <= 0 {
		return models.Invite{}, errx.ErrorInitiatorRoleHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	invID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := e.jwt.CreateInviteToken(jwtmanager.InvitePayload{
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

	stmt := dbx.Invite{
		ID:            invID,
		Status:        enum.InviteStatusSent,
		Role:          params.Role,
		DistributorID: initiator.DistributorID,
		ExpiresAt:     exAt,
		CreatedAt:     now,
	}

	err = e.invite.New().Insert(ctx, stmt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return inviteFromDB(stmt, token), nil
}
