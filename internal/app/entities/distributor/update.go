package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
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
	update := map[string]any{}

	if params.Name != nil {
		update["name"] = *params.Name
	}
	if params.Icon != nil {
		update["icon"] = *params.Icon
	}
	update["updated_at"] = time.Now().UTC()

	distributorQ := d.distributor.New().FilterID(distributorID)

	err := distributorQ.Update(ctx, update)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(fmt.Errorf("updating distributor name: %w", err))
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.ErrorDistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{}, errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}
