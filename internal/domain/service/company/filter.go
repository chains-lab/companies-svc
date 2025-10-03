package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
)

type Filters struct {
	Name     *string
	Statuses []string
}

func (s Service) Filter(
	ctx context.Context,
	filters Filters,
	page, size uint64,
) (models.CompanyCollection, error) {
	res, err := s.db.FilterCompanies(ctx, filters, page, size)
	if err != nil {
		return models.CompanyCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter companies, cause: %w", err),
		)
	}

	return res, err
}
