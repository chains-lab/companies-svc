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

func (s Service) UpdateStatus(
	ctx context.Context,
	distributorID uuid.UUID,
	status string,
) (models.Distributor, error) {
	err := enum.CheckDistributorBlockStatus(status)
	if err != nil {
		return models.Distributor{}, errx.ErrorInvalidDistributorBlockStatus.Raise(
			fmt.Errorf("invalid status %s: %w", status, err),
		)
	}

	if status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.ErrorCannotSetDistributorStatusBlocked.Raise(
			fmt.Errorf("cannot set status to blocked"),
		)
	}

	distributor, err := s.Get(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.ErrorDistributorIsBlocked.Raise(
			fmt.Errorf("distributor %s is blocked", distributorID),
		)
	}

	now := time.Now().UTC()
	err = s.db.UpdateDistributorStatus(ctx, distributorID, status, now)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    enum.DistributorStatusInactive,
		UpdatedAt: now,
		CreatedAt: distributor.CreatedAt,
	}, nil
}
