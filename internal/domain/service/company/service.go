package company

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	event EventPublisher
	db    database
}

func NewService(db database, publisher EventPublisher) Service {
	return Service{
		event: publisher,
		db:    db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmployee(ctx context.Context, input models.Employee) error

	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (models.Employee, error)
	GetEmployeeByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (models.Employee, error)
	GetEmployeeByCompanyAndUserAndRole(ctx context.Context, companyID, userID uuid.UUID, role string) (models.Employee, error)

	CreateCompany(ctx context.Context, input models.Company) (models.Company, error)
	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)
	FilterCompanies(ctx context.Context, filters FiltersParams, page, size uint64) (models.CompaniesCollection, error)

	GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error)

	UpdateCompany(ctx context.Context, ID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	UpdateCompaniesStatus(ctx context.Context, ID uuid.UUID, status string, updatedAt time.Time) error
	DeleteCompany(ctx context.Context, ID uuid.UUID) error
}

type EventPublisher interface {
	PublishCompanyCreated(
		ctx context.Context,
		company models.Company,
	) error

	PublishCompanyDeleted(
		ctx context.Context,
		company models.Company,
		recipients []uuid.UUID,
	) error

	PublishCompanyDeactivated(
		ctx context.Context,
		company models.Company,
		recipients []uuid.UUID,
	) error
	PublishCompanyActivated(
		ctx context.Context,
		company models.Company,
		recipients []uuid.UUID,
	) error

	PublishEmployeeCreated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients []uuid.UUID,
	) error
}
