package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/entities/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
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
