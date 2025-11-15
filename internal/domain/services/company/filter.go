package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
)

type FiltersParams struct {
	Name     *string
	Statuses []string
}

func (s Service) Filter(
	ctx context.Context,
	filters FiltersParams,
	page, size uint64,
) (models.CompaniesCollection, error) {
	res, err := s.db.FilterCompanies(ctx, filters, page, size)
	if err != nil {
		return models.CompaniesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter companies, cause: %w", err),
		)
	}

	return res, err
}
