package distributor

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) CancelBlock(
	ctx context.Context,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	dis, err := s.GetDistributor(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	_, err = s.GetActiveDistributorBlock(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	canceledAt := time.Now().UTC()

	trErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CancelActiveDistributorBlock(ctx, distributorID, canceledAt)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("cancelling active distributor block: %w", err),
			)
		}

		err = s.db.UpdateDistributorStatus(ctx, distributorID, enum.DistributorStatusInactive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		return nil
	})
	if trErr != nil {
		return models.Distributor{}, trErr
	}

	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating distributor status: %w", err),
		)
	}

	dis.Status = enum.DistributorStatusInactive
	dis.UpdatedAt = canceledAt

	return dis, nil
}
