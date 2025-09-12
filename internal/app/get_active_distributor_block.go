package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.Block, error) {
	return a.distributor.GetActiveDistributorBlock(ctx, distributorID)
}
