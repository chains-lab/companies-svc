package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (d Distributor) GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.Block, error) {
	block, err := d.block.New().FilterDistributorID(distributorID).FilterStatus(enum.DistributorBlockStatusActive).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Block{}, errx.ErrorNoActiveBlockForDistributor.Raise(
				fmt.Errorf("active block for distributor with ID %s not found, cause: %w", distributorID, err),
			)
		default:
			return models.Block{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting active block for distributor with ID %s, cause: %w", distributorID, err),
			)
		}
	}

	res := models.Block{
		ID:            block.ID,
		DistributorID: block.DistributorID,
		InitiatorID:   block.InitiatorID,
		Reason:        block.Reason,
		Status:        block.Status,
		BlockedAt:     block.BlockedAt,
	}
	if block.CanceledAt != nil {
		res.CanceledAt = block.CanceledAt
	}

	return res, nil
}
