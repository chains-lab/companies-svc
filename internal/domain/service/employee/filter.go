package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterParams struct {
	UserID    []uuid.UUID
	CompanyID []uuid.UUID
	Roles     []string
}

func (s Service) Filter(
	ctx context.Context,
	filters FilterParams,
	page, size uint64,
) (models.EmployeesCollection, error) {
	res, err := s.db.FilterEmployees(ctx, filters, page, size)
	if err != nil {
		return models.EmployeesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list employees, cause: %w", err),
		)
	}

	userIDs := make([]uuid.UUID, 0, len(res.Data))
	for _, emp := range res.Data {
		userIDs = append(userIDs, emp.UserID)
	}

	return res, nil
}
