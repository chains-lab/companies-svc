package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetBlock(
	ctx context.Context,
	BlockID uuid.UUID,
) (models.DistributorBlock, error) {
	block, err := s.db.GetDistributorBlockByID(ctx, BlockID)
	if err != nil {
		return models.DistributorBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting block with ID %s: %w", BlockID, err),
		)
	}

	if block.IsNil() {
		return models.DistributorBlock{}, errx.ErrorDistributorBlockNotFound.Raise(
			fmt.Errorf("block with ID %s not found", BlockID),
		)
	}

	return block, nil
}

func (s Service) GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.DistributorBlock, error) {
	block, err := s.db.GetActiveDistributorBlock(ctx, distributorID)
	if err != nil {
		return models.DistributorBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting active block for distributor with ID %s, cause: %w", distributorID, err),
		)
	}

	if block.IsNil() {
		return models.DistributorBlock{}, errx.ErrorDistributorBlockNotFound.Raise(
			fmt.Errorf("active block for distributor with ID %s not found", distributorID),
		)
	}

	return block, nil
}
