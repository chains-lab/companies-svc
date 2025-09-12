package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/entities/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FilterBlockagesList struct {
	Distributors []uuid.UUID
	Initiators   []uuid.UUID
	Statuses     []string
}

func (a App) ListBlockages(
	ctx context.Context,
	filter FilterBlockagesList,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Block, pagi.Response, error) {
	input := distributor.FilterBlockagesList{}
	if filter.Distributors != nil && len(filter.Distributors) > 0 {
		input.Distributors = filter.Distributors
	}
	if filter.Initiators != nil && len(filter.Initiators) > 0 {
		input.Initiators = filter.Initiators
	}
	if filter.Statuses != nil && len(filter.Statuses) > 0 {
		input.Statuses = filter.Statuses
	}

	return a.distributor.ListBlockages(ctx, input, pag, sort)
}
