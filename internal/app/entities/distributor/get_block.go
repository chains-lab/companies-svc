package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
)

func (d Distributor) GetBlock(
	ctx context.Context,
	ID uuid.UUID,
) (models.Block, error) {
	block, err := d.block.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Block{}, errx.DistributorBlockNotFound.Raise(
				fmt.Errorf("block with ID %s not found: %w", ID, err),
			)
		default:
			return models.Block{}, errx.ErrorInternal.Raise(
				fmt.Errorf("getting block with ID %s: %w", ID, err),
			)
		}
	}

	return models.Block{
		ID:            block.ID,
		DistributorID: block.DistributorID,
		InitiatorID:   block.InitiatorID,
		Reason:        block.Reason,
		Status:        block.Status,
		BlockedAt:     block.BlockedAt,
		CanceledAt:    block.CanceledAt,
	}, nil
}
