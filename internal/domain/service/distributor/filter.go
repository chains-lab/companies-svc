package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
)

type FilterList struct {
	Name     *string
	Statuses []string
}

func (s Service) List(
	ctx context.Context,
	filters FilterList,
	page, size uint64,
) (models.DistributorCollection, error) {
	res, err := s.db.FilterDistributors(ctx, filters, page, size)
	if err != nil {
		return models.DistributorCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter distributors, cause: %w", err),
		)
	}

	return res, err
}
