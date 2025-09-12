package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.Employee, error) {
	return a.employee.GetInitiator(ctx, initiatorID)
}
