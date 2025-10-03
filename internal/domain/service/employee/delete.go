package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, initiatorID, userID, distributorID uuid.UUID) error {
	_, err := s.Get(ctx, GetFilters{
		UserID:        &userID,
		DistributorID: &distributorID,
	})
	if err != nil {
		return err
	}

	//TODO compare roles of initiator and user

	err = s.db.DeleteEmployee(ctx, userID, distributorID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete employee, cause: %w", err),
		)
	}

	return nil
}
