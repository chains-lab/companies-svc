package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) CreateDistributor(ctx context.Context, initiatorID uuid.UUID, name, icon string) (models.Distributor, error) {
	_, err := a.GetInitiator(ctx, initiatorID)
	if err != nil && !errors.Is(err, errx.ErrorInitiatorNotEmployee) {
		return models.Distributor{}, err
	}
	if err == nil {
		return models.Distributor{}, errx.ErrorCurrentEmployeeCannotCreateDistributor.Raise(
			fmt.Errorf("current employee %s can not create distributor", initiatorID),
		)
	}

	var dist models.Distributor
	trErr := a.transaction(func(ctx context.Context) error {
		dist, err = a.distributor.Create(ctx, distributor.CreateParams{
			Name: name,
			Icon: icon,
		})
		if err != nil {
			return err
		}

		_, err = a.employee.CreateEmployee(ctx, employee.CreateParams{
			UserID:        initiatorID,
			DistributorID: dist.ID,
			Role:          enum.EmployeeRoleOwner,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if trErr != nil {
		return models.Distributor{}, trErr
	}

	return dist, nil
}

func (a App) GetDistributor(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error) {
	return a.distributor.GetDistributor(ctx, distributorID)
}

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
