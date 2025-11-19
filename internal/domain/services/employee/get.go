package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, companyID, userID uuid.UUID) (models.Employee, error) {
	employee, err := s.db.GetEmployee(ctx, userID, companyID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee user %s in company %s, cause: %w", userID, companyID, err),
		)
	}
	if employee.IsNil() {
		return models.Employee{}, errx.ErrorEmployeeNotFound.Raise(
			fmt.Errorf("employee for user EmployeeID %s not found in company %s", userID, companyID),
		)
	}

	return employee, nil
}
