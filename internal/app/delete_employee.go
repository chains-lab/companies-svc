package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) DeleteEmployee(ctx context.Context, initiatorID, employeeID, distributorID uuid.UUID) error {
	_, err := a.GetDistributor(ctx, distributorID)
	if err != nil {
		return err
	}

	return a.employee.Delete(ctx, initiatorID, employeeID, distributorID)
}
