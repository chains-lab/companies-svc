package data

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateEmployee(ctx context.Context, input models.Employee) error {
	schema := employeeModelToSchema(input)

	return d.sql.employees.New().Insert(ctx, schema)
}

func (d *Database) FilterEmployees(
	ctx context.Context,
	filter employee.Filter,
	page, size uint64,
) (models.EmployeeCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.employees.New()

	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}
	if filter.Roles != nil && len(filter.Roles) > 0 {
		query = query.FilterRole(filter.Roles...)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.EmployeeCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.EmployeeCollection{}, err
	}

	collection := make([]models.Employee, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, employeeSchemaToModel(row))
	}

	return models.EmployeeCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) GetEmployee(ctx context.Context, filter employee.GetFilters) (models.Employee, error) {
	query := d.sql.employees.New()

	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}
	if filter.UserID != nil {
		query = query.FilterUserID(*filter.UserID)
	}
	if filter.Role != nil {
		query = query.FilterRole(*filter.Role)
	}

	row, err := query.Get(ctx)
	if err != nil {
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) GetUserEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	row, err := d.sql.employees.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		return models.Employee{}, err
	}

	return employeeSchemaToModel(row), nil
}

func (d *Database) UpdateEmployeeRole(ctx context.Context, userID uuid.UUID, newRole string, updatedAt time.Time) error {
	return d.sql.employees.New().FilterUserID(userID).UpdateRole(newRole).Update(ctx, updatedAt)
}

func (d *Database) DeleteEmployee(ctx context.Context, userID, distributorID uuid.UUID) error {
	return d.sql.employees.New().FilterUserID(userID).FilterDistributorID(distributorID).Delete(ctx)
}

func employeeModelToSchema(input models.Employee) pgdb.Employee {
	return pgdb.Employee{
		UserID:        input.UserID,
		DistributorID: input.DistributorID,
		Role:          input.Role,
		UpdatedAt:     input.UpdatedAt,
		CreatedAt:     input.CreatedAt,
	}
}

func employeeSchemaToModel(input pgdb.Employee) models.Employee {
	return models.Employee{
		UserID:        input.UserID,
		DistributorID: input.DistributorID,
		Role:          input.Role,
		UpdatedAt:     input.UpdatedAt,
		CreatedAt:     input.CreatedAt,
	}
}
