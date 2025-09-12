package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FilterListEmployee struct {
	Distributors []uuid.UUID
	Roles        []string
}

func (e Employee) ListEmployees(
	ctx context.Context,
	filters FilterListEmployee,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Employee, pagi.Response, error) {
	query := e.employee.New()

	if filters.Distributors != nil {
		query = query.FilterDistributorID(filters.Distributors...)
	}
	if filters.Roles != nil {
		for _, role := range filters.Roles {
			err := enum.CheckEmployeeRole(role)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidEmployeeRole.Raise(err)
			}
		}
		query = query.FilterRole(filters.Roles...)
	}

	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	query = query.Page(limit, offset)

	for _, s := range sort {
		ascend := s.Ascend
		switch s.Field {
		case "role":
			query = query.OrderByRole(ascend)
		default:

		}
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	var result []models.Employee
	for _, emp := range rows {
		result = append(result, models.Employee{
			UserID:        emp.UserID,
			DistributorID: emp.DistributorID,
			Role:          emp.Role,
			UpdatedAt:     emp.UpdatedAt,
			CreatedAt:     emp.CreatedAt,
		})
	}

	return result, pagi.Response{
		Total: total,
		Page:  pag.Page,
		Size:  pag.Size,
	}, nil
}
