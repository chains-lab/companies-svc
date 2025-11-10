package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateEmployee(ctx context.Context, input models.Employee) error {
	return d.sql.employees.New().Insert(ctx, employeeModelToSchema(input))
}

func (d *Database) GetEmployeeByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (models.Employee, error) {
	row, err := d.sql.employees.New().FilterUserID(userID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) GetEmployeeByCompanyAndUser(
	ctx context.Context,
	companyID, userID uuid.UUID,
) (models.Employee, error) {
	row, err := d.sql.employees.New().FilterCompanyID(companyID).FilterUserID(userID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) GetEmployeeByCompanyAndUserAndRole(
	ctx context.Context,
	companyID, userID uuid.UUID,
	role string,
) (models.Employee, error) {
	row, err := d.sql.employees.New().FilterCompanyID(companyID).FilterUserID(userID).FilterRole(role).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) GetEmployee(
	ctx context.Context,
	params employee.GetParams,
) (models.Employee, error) {
	query := d.sql.employees.New()

	if params.UserID != nil {
		query = query.FilterUserID(*params.UserID)
	}
	if params.CompanyID != nil {
		query = query.FilterCompanyID(*params.CompanyID)
	}
	if params.Role != nil {
		query = query.FilterRole(*params.Role)
	}

	row, err := query.Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Employee{}, nil
	case err != nil:
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) FilterEmployees(
	ctx context.Context,
	filter employee.FilterParams,
	page, size uint64,
) (models.EmployeesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.employees.New()

	if filter.CompanyID != nil {
		query = query.FilterCompanyID(filter.CompanyID...)
	}
	if filter.UserID != nil {
		query = query.FilterUserID(filter.UserID...)
	}
	if filter.Roles != nil && len(filter.Roles) > 0 {
		query = query.FilterRole(filter.Roles...)
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

func (d *Database) UpdateEmployee(
	ctx context.Context,
	userID uuid.UUID,
	params employee.UpdateEmployeeParams,
	updatedAt time.Time,
) error {
	q := d.sql.employees.New().FilterUserID(userID)
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

	if empty {
		return nil
	}

	return q.Update(ctx, updatedAt)
}

func (d *Database) UpdateEmployeeRole(
	ctx context.Context,
	userID uuid.UUID,
	role string,
	updatedAt time.Time,
) error {
	return d.sql.employees.New().FilterUserID(userID).UpdateRole(role).Update(ctx, updatedAt)
}

func (d *Database) DeleteEmployee(ctx context.Context, userID, companyID uuid.UUID) error {
	return d.sql.employees.New().FilterUserID(userID).FilterCompanyID(companyID).Delete(ctx)
}

func (d *Database) EmployeeExist(ctx context.Context, userID uuid.UUID) (bool, error) {
	return d.sql.employees.New().FilterUserID(userID).Exist(ctx)
}

func (d *Database) GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error) {
	query := d.sql.employees.New().FilterCompanyID(companyID)
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
