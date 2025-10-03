package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterBlockages struct {
	DistributorID *uuid.UUID
	InitiatorID   *uuid.UUID
	Status        *string
}

func (s Service) FilterBlockages(
	ctx context.Context,
	filters FilterBlockages,
	page, size uint64,
) (models.DistributorBlockCollection, error) {
	res, err := s.db.FilterDistributorBlocks(ctx, filters, page, size)
	if err != nil {
		return models.DistributorBlockCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter distributor blocks, cause: %w", err),
		)
	}

	return res, err
}
