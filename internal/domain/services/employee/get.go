package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type GetParams struct {
	UserID    *uuid.UUID
	CompanyID *uuid.UUID
	Role      *string
}

func (s Service) Get(ctx context.Context, params GetParams) (models.Employee, error) {
	res, err := s.db.GetEmployee(ctx, params)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee, cause: %w", err),
		)
	}
	if res.IsNil() {
		return models.Employee{}, errx.ErrorEmployeeNotFound.Raise(
			fmt.Errorf("employee not found"),
		)
	}

	return res, nil
}

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.Employee, error) {
	employee, err := s.db.GetEmployeeByUserID(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user ID %s, cause: %w", initiatorID, err),
		)
	}
	if employee.IsNil() {
		return models.Employee{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("employee not found for user %s", initiatorID),
		)
	}

	return employee, nil
}
