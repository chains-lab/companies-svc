package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (d Distributor) Block(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.Block, error) {
	blockages := dbx.Blockages{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Status:        enum.DistributorBlockStatusActive,
		BlockedAt:     time.Now().UTC(),
	}

	_, err := d.block.New().FilterDistributorID(distributorID).FilterStatus(enum.DistributorBlockStatusActive).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Block{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if err == nil {
		return models.Block{}, errx.DistributorHaveAlreadyActiveBlock.Raise(
			fmt.Errorf("distributor %s already has an active block", distributorID),
		)
	}

	trErr := d.distributor.Transaction(func(ctx context.Context) error {
		err = d.distributor.New().FilterID(distributorID).Update(ctx, map[string]any{
			"status":     enum.DistributorStatusBlocked,
			"updated_at": time.Now().UTC(),
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		_, err = d.block.New().FilterDistributorID(distributorID).FilterStatus(enum.DistributorBlockStatusActive).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("checking existing block: %w", err),
			)
		}
		if err == nil {
			return errx.DistributorHaveAlreadyActiveBlock.Raise(
				fmt.Errorf("distributor %s already has an active block", distributorID),
			)
		}

		err = d.block.Insert(ctx, blockages)
		if err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("inserting new block: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.Block{}, trErr
	}

	return models.Block{
		ID:            blockages.ID,
		DistributorID: blockages.DistributorID,
		InitiatorID:   blockages.InitiatorID,
		Reason:        blockages.Reason,
		BlockedAt:     blockages.BlockedAt,
	}, nil
}
