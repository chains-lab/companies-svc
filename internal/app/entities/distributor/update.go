package distributor

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Name *string
	Icon *string
}

func (d Distributor) Update(ctx context.Context,
	distributorID uuid.UUID,
	params UpdateParams,
) (models.Distributor, error) {
	distributor, err := d.GetDistributor(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.ErrorDistributorIsBlocked.Raise(
			fmt.Errorf("distributor %s is blocked", distributorID),
		)
	}

	update := map[string]any{}

	if params.Name != nil {
		update["name"] = *params.Name
		distributor.Name = *params.Name
	}
	if params.Icon != nil {
		update["icon"] = *params.Icon
		distributor.Icon = *params.Icon
	}
	distributor.UpdatedAt = time.Now().UTC()
	update["updated_at"] = distributor.UpdatedAt

	err = d.distributor.New().FilterID(distributorID).Update(ctx, update)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(fmt.Errorf("updating distributor name: %w", err))
	}

	return distributor, nil
}
