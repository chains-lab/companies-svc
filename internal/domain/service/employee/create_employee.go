package employee

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type CreateParams struct {
	UserID        uuid.UUID
	DistributorID uuid.UUID
	Role          string
}

func (s Service) Create(ctx context.Context, params CreateParams) (models.Employee, error) {
	emp, err := s.db.GetEmployee(ctx, GetFilters{
		UserID: &params.UserID,
	})
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("check existing employee: %w", err),
		)
	}

	if !emp.IsNil() {
		return models.Employee{}, errx.ErrorEmployeeAlreadyExists.Raise(
			fmt.Errorf("employee with user ID %s already exists", params.UserID),
		)
	}

	now := time.Now().UTC()
	err = enum.CheckEmployeeRole(params.Role)
	if err != nil {
		return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
			errors.New("invalid employee role"),
		)
	}

	emp = models.Employee{
		UserID:        params.UserID,
		DistributorID: params.DistributorID,
		Role:          params.Role,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.db.CreateEmployee(ctx, emp)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create employee, cause: %w", err),
		)
	}

	return emp, nil
}
