package company

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
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
) (models.Company, error) {
	now := time.Now().UTC()

	dis := models.Company{
		ID:        uuid.New(),
		Name:      params.Name,
		Icon:      params.Icon,
		Status:    enum.DistributorStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := s.db.GetUserEmployee(ctx, params.InitiatorID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user employee: %w", err),
		)
	}
	if !user.IsNil() {
		return models.Company{}, errx.ErrorCurrentEmployeeCannotCreatecompany.Raise(
			fmt.Errorf("user is already an employee of company: %s", user.CompanyID),
		)
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		_, err = s.db.CreateCompany(ctx, dis)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}

		emp := models.Employee{
			UserID:    params.InitiatorID,
			CompanyID: dis.ID,
			Role:      enum.EmployeeRoleAdmin,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = s.db.CreateEmployee(ctx, emp)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create employee for company: %w", err),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Company{}, txErr
	}

	return dis, err
}
