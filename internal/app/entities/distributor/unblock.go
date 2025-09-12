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

func (d Distributor) Unblock(
	ctx context.Context,
	blockID uuid.UUID,
) (models.Block, error) {
	block, err := d.GetBlock(ctx, blockID)

	canceledAt := time.Now().UTC()
	trErr := d.distributor.Transaction(func(ctx context.Context) error {
		err = d.block.FilterID(blockID).FilterStatus(enum.BlockStatusActive).Update(ctx, map[string]any{
			"status":       enum.BlockStatusCanceled,
			"cancelled_at": canceledAt,
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating block status: %w", err),
			)
		}

		err = d.distributor.New().FilterID(block.DistributorID).Update(ctx, map[string]any{
			"status":     enum.DistributorStatusInactive,
			"updated_at": canceledAt,
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		return nil
	})
	if trErr != nil {
		return models.Block{}, trErr
	}

	return models.Block{
		ID:            block.ID,
		DistributorID: block.DistributorID,
		InitiatorID:   block.InitiatorID,
		Reason:        block.Reason,
		Status:        enum.BlockStatusCanceled,
		BlockedAt:     block.BlockedAt,
		CanceledAt:    &canceledAt,
	}, nil
}
