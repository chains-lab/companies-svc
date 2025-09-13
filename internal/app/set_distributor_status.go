package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (a App) SetDistributorStatus(ctx context.Context, initiatorID, distributorID uuid.UUID, status string) (models.Distributor, error) {
	initiator, err := a.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Distributor{}, err
	}
	if initiator.Role != enum.EmployeeRoleOwner {
		return models.Distributor{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("current employee %s can not set distributor status onlt owner can did it", initiatorID),
		)
	}

	var res models.Distributor

	switch status {
	case enum.DistributorStatusInactive:
		trErr := a.transaction(func(ctx context.Context) error {
			res, err = a.distributor.SetStatusInactive(ctx, distributorID)
			if err != nil {
				return err
			}

			err = a.employee.DeleteMany(ctx, employee.DeleteManyParams{
				DistributorID: &distributorID,
				Roles:         []string{enum.EmployeeRoleAdmin, enum.EmployeeRoleModerator},
			})
			if err != nil {
				return err
			}

			return nil
		})
		if trErr != nil {
			return models.Distributor{}, trErr
		}

		//TODO event in kafka for services

	case enum.DistributorStatusActive:
		res, err = a.distributor.SetStatusActive(ctx, distributorID)
		if err != nil {
			return models.Distributor{}, err
		}

	default:
		return models.Distributor{}, errx.ErrorUnexpectedDistributorSetStatus.Raise(
			fmt.Errorf("unexpected status: %s", status),
		)
	}

	return res, err
}
