package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterBlockagesList struct {
	Distributors []uuid.UUID
	Initiators   []uuid.UUID
	Statuses     []string
}

func (s Service) ListBlockages(
	ctx context.Context,
	filters FilterBlockagesList,
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
