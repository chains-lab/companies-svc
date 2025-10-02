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

type CreateParams struct {
	Name string
	Icon string
}

func (s Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.Distributor, error) {
	now := time.Now().UTC()

	dis := models.Distributor{
		ID:        uuid.New(),
		Name:      params.Name,
		Icon:      params.Icon,
		Status:    enum.DistributorStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.db.CreateDistributor(ctx, dis)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return dis, err
}
