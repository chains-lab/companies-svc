package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (r *Repo) CreateEmployee(ctx context.Context, input models.Employee) error {
	return r.sql.employees.New().Insert(ctx, employeeModelToSchema(input))
}

func (r *Repo) GetEmployee(
	ctx context.Context,
	ID uuid.UUID,
) (models.Employee, error) {
	row, err := r.sql.employees.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (r *Repo) GetEmployeeUserInCompany(
	ctx context.Context,
	userID, companyID uuid.UUID,
) (models.Employee, error) {
	row, err := r.sql.employees.New().FilterCompanyID(companyID).FilterUserID(userID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (r *Repo) GetCompanyOwner(ctx context.Context, companyID uuid.UUID) (models.Employee, error) {
	row, err := r.sql.employees.New().FilterCompanyID(companyID).FilterRole(enum.EmployeeRoleOwner).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (r *Repo) FilterEmployees(
	ctx context.Context,
	filter employee.FilterParams,
	page, size uint64,
) (models.EmployeesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := r.sql.employees.New()

	if filter.CompanyID != nil {
		query = query.FilterCompanyID(filter.CompanyID...)
	}
	if filter.UserID != nil {
		query = query.FilterUserID(filter.UserID...)
	}
	if filter.Roles != nil && len(filter.Roles) > 0 {
		query = query.FilterRole(filter.Roles...)
	}
	if filter.EmployeeID != nil && len(filter.EmployeeID) > 0 {
		query = query.FilterID(filter.EmployeeID...)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.EmployeesCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.EmployeesCollection{}, err
	}

	collection := make([]models.Employee, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, employeeSchemaToModel(row))
	}

	return models.EmployeesCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (r *Repo) UpdateEmployee(
	ctx context.Context,
	ID uuid.UUID,
	params employee.UpdateParams,
	updatedAt time.Time,
) error {
	q := r.sql.employees.New().FilterID(ID)
	empty := true

	if params.Position != nil {
		if *params.Position == "" {
			q = q.UpdatePosition(nil)
		} else {
			q = q.UpdatePosition(params.Position)
		}
		empty = false
	}

	if params.Label != nil {
		if *params.Label == "" {
			q = q.UpdateLabel(nil)
		} else {
			q = q.UpdateLabel(params.Label)
		}
		empty = false
	}

	if params.Role != nil {
		q = q.UpdateRole(*params.Role)
		empty = false
	}

	if empty {
		return nil
	}

	return q.Update(ctx, updatedAt)
}

func (r *Repo) DeleteEmployee(ctx context.Context, ID uuid.UUID) error {
	return r.sql.employees.New().FilterID(ID).Delete(ctx)
}

func (r *Repo) GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error) {
	query := r.sql.employees.New().FilterCompanyID(companyID)
	if len(roles) > 0 {
		query = query.FilterRole(roles...)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return models.EmployeesCollection{}, err
	}

	collection := make([]models.Employee, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, employeeSchemaToModel(row))
	}

	return models.EmployeesCollection{
		Data:  collection,
		Page:  1,
		Size:  uint64(len(collection)),
		Total: uint64(len(collection)),
	}, nil
}
