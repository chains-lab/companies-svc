package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (d Distributor) Unblock(
	ctx context.Context,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	dis, err := d.GetDistributor(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	canceledAt := time.Now().UTC()

	err = d.block.FilterDistributorID(distributorID).FilterStatus(enum.DistributorBlockStatusActive).Update(ctx, map[string]any{
		"status":       enum.DistributorBlockStatusCanceled,
		"cancelled_at": canceledAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.ErrorNoActiveBlockForDistributor.Raise(
				fmt.Errorf("no active block for distributor %s: %w", distributorID, err),
			)
		default:
			return models.Distributor{}, errx.ErrorInternal.Raise(
				fmt.Errorf("updating block status: %w", err),
			)
		}
	}

	err = d.distributor.New().FilterID(distributorID).Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": canceledAt,
	})
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating distributor status: %w", err),
		)
	}

	return d.GetDistributor(ctx, dis.ID)
}
