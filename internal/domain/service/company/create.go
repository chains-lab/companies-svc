package company

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type CreateParams struct {
	Name string
	Icon string
}

func (s Service) Create(
	ctx context.Context,
	initiatorID uuid.UUID,
	params CreateParams,
) (models.Company, error) {
	now := time.Now().UTC()

	comp := models.Company{
		ID:        uuid.New(),
		Name:      params.Name,
		Icon:      params.Icon,
		Status:    enum.CompanyStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := s.db.GetEmployeeByUserID(ctx, initiatorID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get user employee, cause: %w", err),
		)
	}
	if !user.IsNil() {
		return models.Company{}, errx.ErrorCurrentEmployeeCannotCreateCompany.Raise(
			fmt.Errorf("user is already an employee of company: %s", user.CompanyID),
		)
	}

	var employee models.Employee
	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		_, err = s.db.CreateCompany(ctx, comp)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create company, cause: %w", err),
			)
		}

		employee = models.Employee{
			UserID:    initiatorID,
			CompanyID: comp.ID,
			Role:      enum.EmployeeRoleOwner,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = s.db.CreateEmployee(ctx, employee)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create employee for company, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Company{}, err
	}

	if err = s.event.PublishEmployeeCreated(ctx, employee); err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee with kafka, cause: %w", err),
		)
	}

	if err = s.event.PublishCompanyCreated(ctx, comp); err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company created event, cause: %w", err),
		)
	}

	return comp, err
}
