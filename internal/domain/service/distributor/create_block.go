package distributor

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) CreteBlock(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.DistributorBlock, error) {
	_, err := s.GetDistributor(ctx, distributorID)
	if err != nil {
		return models.DistributorBlock{}, err
	}

	block := models.DistributorBlock{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Status:        enum.DistributorBlockStatusActive,
		BlockedAt:     time.Now().UTC(),
	}

	activeBlock, err := s.db.GetActiveDistributorBlock(ctx, block.ID)
	if err != nil {
		return models.DistributorBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if !activeBlock.IsNil() {
		return models.DistributorBlock{}, errx.DistributorHaveAlreadyActiveBlock.Raise(
			fmt.Errorf("distributor %s already has an active block", distributorID),
		)
	}

	trErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.UpdateDistributorStatus(ctx, distributorID, enum.DistributorStatusBlocked)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		err = s.db.CreateDistributorBlock(ctx, block)
		if err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("inserting new block: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.DistributorBlock{}, trErr
	}

	return block, nil
}
