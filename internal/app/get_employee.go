package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	return a.employee.GetByUserID(ctx, userID)
}
