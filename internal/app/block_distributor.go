package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) BlockDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.Block, error) {
	return a.distributor.Block(ctx, initiatorID, distributorID, reason)
}
