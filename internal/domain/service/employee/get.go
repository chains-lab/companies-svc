package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type GetFilters struct {
	UserID    *uuid.UUID
	CompanyID *uuid.UUID
	Role      *string
}

func (s Service) Get(ctx context.Context, filters GetFilters) (models.Employee, error) {
	employee, err := s.db.GetEmployee(ctx, filters)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get employee by user ID, cause: %w", err),
		)
	}

	if employee.IsNil() {
		return models.Employee{}, errx.ErrorEmployeeNotFound.Raise(
			fmt.Errorf("employee not found"),
		)
	}

	return employee, nil
}

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.Employee, error) {
	employee, err := s.db.GetEmployee(ctx, GetFilters{
		UserID: &initiatorID,
	})
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get employee by user ID, cause: %w", err),
		)
	}

	if employee.IsNil() {
		return models.Employee{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("employee not found"),
		)
	}

	return employee, nil
}
