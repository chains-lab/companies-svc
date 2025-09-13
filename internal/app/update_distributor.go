package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/entities/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type UpdateDistributorParams struct {
	Name *string
	Icon *string
}

func (a App) UpdateDistributor(ctx context.Context, initiatorID, distributorID uuid.UUID, params UpdateDistributorParams) (models.Distributor, error) {
	initiator, err := a.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if initiator.Role != enum.EmployeeRoleAdmin && initiator.Role != enum.EmployeeRoleOwner {
		return models.Distributor{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator role has not enough rights only %s or %s can update distributor, got %s",
				enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner, initiator.Role,
			),
		)
	}

	input := distributor.UpdateParams{}
	if params.Name != nil {
		input.Name = params.Name
	}
	if params.Icon != nil {
		input.Icon = params.Icon
	}

	return a.distributor.Update(ctx, distributorID, input)
}
