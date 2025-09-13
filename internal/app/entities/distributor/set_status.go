package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (d Distributor) SetStatusInactive(
	ctx context.Context,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	distributor, err := d.GetDistributor(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.ErrorDistributorIsBlocked.Raise(
			fmt.Errorf("distributor %s is blocked", distributorID),
		)
	}

	distributorQ := d.distributor.New().FilterID(distributorID)

	now := time.Now().UTC()
	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": now,
	})
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

func (d Distributor) SetStatusActive(
	ctx context.Context,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	distributor, err := d.distributor.New().FilterID(distributorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.ErrorDistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{},
				errx.ErrorInternal.Raise(
					fmt.Errorf("internal error: %w", err),
				)
		}
	}
	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.ErrorDistributorIsBlocked.Raise(
			fmt.Errorf("distributor %s is block", distributorID),
		)
	}

	distributorQ := d.distributor.New().FilterID(distributorID)

	now := time.Now().UTC()
	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusActive,
		"updated_at": now,
	})
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(fmt.Errorf("updating distributor status: %w", err))
	}

	distributor, err = distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.ErrorDistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{},
				errx.ErrorInternal.Raise(
					fmt.Errorf("internal error: %w", err),
				)
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    enum.DistributorStatusActive,
		UpdatedAt: now,
		CreatedAt: distributor.CreatedAt,
	}, nil
}
