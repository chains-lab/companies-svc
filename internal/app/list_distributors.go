package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/entities/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/pagi"
)

type FilterDistributorList struct {
	Name     *string
	Statuses []string
}

func (a App) ListDistributors(
	ctx context.Context,
	filter FilterDistributorList,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Distributor, pagi.Response, error) {
	input := distributor.FilterList{}
	if filter.Name != nil {
		input.Name = filter.Name
	}
	if filter.Statuses != nil && len(filter.Statuses) > 0 {
		input.Statuses = filter.Statuses
	}
	return a.distributor.List(ctx, input, pag, sort)
}
