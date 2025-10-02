package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type CreateInviteParams struct {
	InitiatorID   uuid.UUID
	DistributorID uuid.UUID
	Role          string
}

func (a App) CreateInvite(ctx context.Context, params CreateInviteParams) (models.Invite, error) {
	distributor, err := a.distributor.GetDistributor(ctx, params.DistributorID)
	if err != nil {
		return models.Invite{}, err
	}

	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Invite{}, errx.ErrorDistributorIsBlocked.Raise(
			fmt.Errorf("distributor %s is blocked", params.DistributorID),
		)
	}

	return a.employee.CreateInvite(ctx, employee.SentInviteParams{
		InitiatorID:   params.InitiatorID,
		DistributorID: params.DistributorID,
		Role:          params.Role,
	})
}

func (a App) AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.employee.AcceptInvite(ctx, userID, token)
		if err != nil {
			return err
		}

		dist, err := a.GetDistributor(ctx, invite.DistributorID)
		if err != nil {
			return err
		}

		if dist.Status != enum.DistributorStatusActive {
			return errx.ErrorAnswerToInviteForNotActiveDistributor.Raise(
				fmt.Errorf("cannot answer to invite for not active distributor %s", dist.ID),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, nil
}
