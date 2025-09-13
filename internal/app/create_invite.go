package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
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
