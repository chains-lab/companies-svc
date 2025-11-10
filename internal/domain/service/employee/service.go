package employee

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

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	CreateEmployee(ctx context.Context, input models.Employee) error

	GetEmployee(ctx context.Context, params GetParams) (models.Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (models.Employee, error)
	GetEmployeeByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (models.Employee, error)
	GetEmployeeByCompanyAndUserAndRole(ctx context.Context, companyID, userID uuid.UUID, role string) (models.Employee, error)
	FilterEmployees(ctx context.Context, filters FilterParams, page, size uint64) (models.EmployeesCollection, error)

	GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error)

	UpdateEmployee(ctx context.Context, userID uuid.UUID, params UpdateEmployeeParams, updatedAt time.Time) error
	UpdateEmployeeRole(ctx context.Context, userID uuid.UUID, role string, updatedAt time.Time) error
	DeleteEmployee(ctx context.Context, userID, companyID uuid.UUID) error
}

type EventPublisher interface {
	PublishEmployeeCreated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients []uuid.UUID,
	) error

	PublishEmployeeUpdated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients []uuid.UUID,
	) error

	PublishEmployeeDeleted(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		Recipients []uuid.UUID,
	) error
}
