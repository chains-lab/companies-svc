package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) RefuseOwnEmployee(ctx context.Context, initiatorID uuid.UUID) error {
	return a.employee.RefuseOwn(ctx, initiatorID)
}
