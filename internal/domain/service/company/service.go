package company

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type EventWriter interface {
	UpdateEmployee(
		ctx context.Context,
		userID uuid.UUID,
		companyID *uuid.UUID,
		role *string,
	) error
}

type Service struct {
	eve EventWriter
	db  database
}

func NewService(db database, eve EventWriter) Service {
	return Service{
		eve: eve,
		db:  db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmployee(ctx context.Context, input models.Employee) error

	GetEmployeeByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (models.Employee, error)

	GetEmployeeByCompanyAndUser(
		ctx context.Context,
		companyID, userID uuid.UUID,
	) (models.Employee, error)

	GetEmployeeByCompanyAndUserAndRole(
		ctx context.Context,
		companyID, userID uuid.UUID,
		role string,
	) (models.Employee, error)

	CreateCompany(ctx context.Context, input models.Company) (models.Company, error)

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	FilterCompanies(
		ctx context.Context,
		filters FiltersParams,
		page, size uint64,
	) (models.CompanyCollection, error)

	UpdateCompany(ctx context.Context, ID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	UpdateCompaniesStatus(ctx context.Context, ID uuid.UUID, status string, updatedAt time.Time) error
}
