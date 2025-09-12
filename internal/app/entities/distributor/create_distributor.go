package distributor

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type CreateParams struct {
	Name string
	Icon string
}

func (d Distributor) Create(
	ctx context.Context,
	params CreateParams,
) (models.Distributor, error) {
	stmt := dbx.Distributor{
		ID:        uuid.New(),
		Name:      params.Name,
		Icon:      params.Icon,
		Status:    enum.DistributorStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := d.distributor.New().Insert(ctx, stmt)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return models.Distributor{
		ID:        stmt.ID,
		Name:      stmt.Name,
		Icon:      stmt.Icon,
		Status:    stmt.Status,
		CreatedAt: stmt.CreatedAt,
		UpdatedAt: stmt.UpdatedAt,
	}, err
}
