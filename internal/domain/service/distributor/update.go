package distributor

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Name *string
	Icon *string
}

func (s Service) Update(ctx context.Context,
	distributorID uuid.UUID,
	params UpdateParams,
) (models.Distributor, error) {
	distributor, err := s.Get(ctx, distributorID)
	if err != nil {
		return models.Distributor{}, err
	}

	mow := time.Now().UTC()

	if params.Name != nil {
		distributor.Name = *params.Name
	}

	if params.Icon != nil {
		distributor.Icon = *params.Icon
	}

	err = s.db.UpdateDistributor(ctx, distributorID, params, mow)
	if err != nil {
		return models.Distributor{}, err
	}

	return distributor, nil
}
