package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Filter struct {
	DistributorID *uuid.UUID
	Roles         []string
}

func (s Service) Filter(
	ctx context.Context,
	filters Filter,
	page, size uint64,
) (models.EmployeeCollection, error) {
	res, err := s.db.FilterEmployees(ctx, filters, page, size)
	if err != nil {
		return models.EmployeeCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list employees, cause: %w", err),
		)
	}

	return res, nil
}
