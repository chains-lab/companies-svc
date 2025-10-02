package employee

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
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

func (s Service) CreateEmployee(ctx context.Context, params CreateParams) (models.Employee, error) {
	_, err := s.GetByUserID(ctx, params.UserID)
	if err != nil && !errors.Is(err, errx.ErrorEmployeeNotFound) {
		return models.Employee{}, err
	}
	if err == nil {
		return models.Employee{}, errx.ErrorEmployeeAlreadyExists.Raise(
			errors.New("employee already exists"),
		)
	}

	now := time.Now().UTC()
	err = enum.CheckEmployeeRole(params.Role)
	if err != nil {
		return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
			errors.New("invalid employee role"),
		)
	}

	err = s.employee.New().Insert(ctx, pgdb.Employee{
		UserID:        params.UserID,
		DistributorID: params.DistributorID,
		Role:          params.Role,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create employee, cause: %w", err),
		)
	}

	return models.Employee{
		UserID:        params.UserID,
		DistributorID: params.DistributorID,
		Role:          params.Role,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
