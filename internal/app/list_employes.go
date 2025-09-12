package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FilterEmployeeList struct {
	Distributors []uuid.UUID
	Roles        []string
}

func (a App) ListEmployees(
	ctx context.Context,
	filter FilterEmployeeList,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Employee, pagi.Response, error) {
	params := employee.FilterListEmployee{}
	if filter.Distributors != nil && len(filter.Distributors) > 0 {
		params.Distributors = filter.Distributors
	}
	if filter.Roles != nil && len(filter.Roles) > 0 {
		params.Roles = filter.Roles
	}

	return a.employee.ListEmployees(ctx, params, pag, sort)
}
