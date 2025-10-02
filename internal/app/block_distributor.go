package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) BlockDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.DistributorBlock, error) {
	return a.distributor.CreteBlock(ctx, initiatorID, distributorID, reason)
}

func (a App) UnblockDistributor(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error) {
	var dis models.Distributor
	var err error

	trErt := a.transaction(func(ctx context.Context) error {
		dis, err = a.distributor.Unblock(ctx, distributorID)
		if err != nil {
			return err
		}

		return nil
	})
	if trErt != nil {
		return models.Distributor{}, trErt
	}

	//TODO Kafka event maybe? Maybe will be better if user unblock his places by own

	return dis, nil
}

func (a App) GetBlock(ctx context.Context, blockID uuid.UUID) (models.DistributorBlock, error) {
	return a.distributor.GetBlock(ctx, blockID)
}

func (a App) GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.DistributorBlock, error) {
	return a.distributor.GetActiveDistributorBlock(ctx, distributorID)
}

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
) ([]models.DistributorBlock, pagi.Response, error) {
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
