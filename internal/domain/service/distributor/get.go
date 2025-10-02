package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetDistributor(ctx context.Context, ID uuid.UUID) (models.Distributor, error) {
	distributor, err := s.db.GetDistributorByID(ctx, ID)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if distributor.IsNil() {
		return models.Distributor{}, errx.ErrorDistributorNotFound.Raise(
			fmt.Errorf("distributor with ID %s not found", ID),
		)
	}

	return distributor, nil
}
