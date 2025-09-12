package distributor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
)

func (d Distributor) GetByID(ctx context.Context, ID uuid.UUID) (models.Distributor, error) {
	distributor, err := d.distributor.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.ErrorDistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", ID, err),
			)
		default:
			return models.Distributor{}, errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Icon:      distributor.Icon,
		Name:      distributor.Name,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}
