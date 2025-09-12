package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetBlock(ctx context.Context, blockID uuid.UUID) (models.Block, error) {
	return a.distributor.GetBlock(ctx, blockID)
}
