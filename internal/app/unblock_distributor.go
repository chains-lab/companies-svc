package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) UnblockDistributor(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error) {
	var dis models.Distributor
	var err error

	trErt := a.transaction(func(ctx context.Context) error {
		dis, err = a.distributor.Unblock(ctx, distributorID)
		if err != nil {
			return err
		}

		return nil
	})
	if trErt != nil {
		return models.Distributor{}, trErt
	}

	//TODO Kafka event maybe? Maybe will be better if user unblock his places by own

	return dis, nil
}
