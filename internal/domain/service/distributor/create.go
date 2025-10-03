package distributor

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type CreateParams struct {
	InitiatorID uuid.UUID
	Name        string
	Icon        string
}

func (s Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.Distributor, error) {
	now := time.Now().UTC()

	dis := models.Distributor{
		ID:        uuid.New(),
		Name:      params.Name,
		Icon:      params.Icon,
		Status:    enum.DistributorStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := s.db.GetUserEmployee(ctx, params.InitiatorID)
	if err != nil {
		return models.Distributor{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user employee: %w", err),
		)
	}
	if !user.IsNil() {
		return models.Distributor{}, errx.ErrorCurrentEmployeeCannotCreateDistributor.Raise(
			fmt.Errorf("user is already an employee of distributor: %s", user.DistributorID),
		)
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		_, err = s.db.CreateDistributor(ctx, dis)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}

		emp := models.Employee{
			UserID:        params.InitiatorID,
			DistributorID: dis.ID,
			Role:          enum.EmployeeRoleAdmin,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		err = s.db.CreateEmployee(ctx, emp)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create employee for distributor: %w", err),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Distributor{}, txErr
	}

	return dis, err
}
