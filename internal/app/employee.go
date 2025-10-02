package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.Employee, error) {
	return a.employee.GetInitiator(ctx, initiatorID)
}

func (a App) GetEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	return a.employee.GetByUserID(ctx, userID)
}

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

func (a App) DeleteEmployee(ctx context.Context, initiatorID, employeeID, distributorID uuid.UUID) error {
	_, err := a.GetDistributor(ctx, distributorID)
	if err != nil {
		return err
	}

	return a.employee.Delete(ctx, initiatorID, employeeID, distributorID)
}

func (a App) RefuseOwnEmployee(ctx context.Context, initiatorID uuid.UUID) error {
	return a.employee.RefuseOwn(ctx, initiatorID)
}
