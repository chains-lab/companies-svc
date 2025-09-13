package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

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
