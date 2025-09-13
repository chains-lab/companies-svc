package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetDistributor(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error) {
	return a.distributor.GetDistributor(ctx, distributorID)
}
