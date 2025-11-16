package employee

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
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}

// create deprecated creates a new employee in the specified company.
func (s Service) create(ctx context.Context, params CreateParams) (models.Employee, error) {
	comp, err := s.getCompany(ctx, params.CompanyID)
	if err != nil {
		return models.Employee{}, err
	}
	if comp.Status != enum.CompanyStatusActive {
		return models.Employee{}, errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("cannot add employee to inactive company with ID %s", params.CompanyID),
		)
	}

	emp, err := s.db.GetEmployeeByUserID(ctx, params.UserID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("falied to check existing employee, cause: %w", err),
		)
	}
	if !emp.IsNil() {
		return models.Employee{}, errx.ErrorUserAlreadyEmployee.Raise(
			fmt.Errorf("employee with user ID %s already exists", params.UserID),
		)
	}

	now := time.Now().UTC()
	err = enum.CheckEmployeeRole(params.Role)
	if err != nil {
		return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("failed to check employee role, cause: %w", err),
		)
	}

	emp = models.Employee{
		UserID:    params.UserID,
		CompanyID: params.CompanyID,
		Role:      params.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var employee models.EmployeesCollection
	err = s.db.Transaction(ctx, func(txCtx context.Context) error {
		err = s.db.CreateEmployee(ctx, emp)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create employee, cause: %w", err),
			)
		}

		employee, err = s.db.GetCompanyEmployees(ctx, emp.CompanyID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get company employees, cause: %w", err),
			)
		}

		return nil
	})

	var recipientIDs []uuid.UUID
	for _, emp := range employee.Data {
		recipientIDs = append(recipientIDs, emp.UserID)
	}

	err = s.event.PublishEmployeeCreated(ctx, comp, emp, recipientIDs...)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish employee created event, cause: %w", err),
		)
	}

	return emp, nil
}
