package company

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
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

	GetEmployee(ctx context.Context, userID, companyID uuid.UUID) (models.Employee, error)
	GetCompanyOwner(ctx context.Context, companyID uuid.UUID) (models.Employee, error)

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
		owner models.Employee,
	) error

	PublishCompanyDeleted(
		ctx context.Context,
		company models.Company,
		recipients ...uuid.UUID,
	) error

	PublishCompanyUpdated(
		ctx context.Context,
		company models.Company,
		recipients ...uuid.UUID,
	) error

	PublishCompanyDeactivated(
		ctx context.Context,
		company models.Company,
		recipients ...uuid.UUID,
	) error

	PublishCompanyActivated(
		ctx context.Context,
		company models.Company,
		recipients ...uuid.UUID,
	) error

	PublishEmployeeCreated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients ...uuid.UUID,
	) error
}

func (s Service) validateInitiator(
	ctx context.Context,
	userID uuid.UUID,
	companyID uuid.UUID,
	roles ...string,
) (models.Employee, error) {
	employee, err := s.db.GetEmployee(ctx, userID, companyID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user EmployeeID, cause: %w", err),
		)
	}
	if employee.IsNil() {
		return models.Employee{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("employee for user EmployeeID %s not found", userID),
		)
	}

	if employee.CompanyID != companyID {
		return models.Employee{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator is not an employee of company: %s", companyID),
		)
	}

	if len(roles) > 0 {
		hasRole := false
		for _, role := range roles {
			if employee.Role == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			return models.Employee{}, errx.ErrorNotEnoughRight.Raise(
				fmt.Errorf("initiator have not enough rights in company: %s", companyID),
			)
		}
	}

	return employee, nil
}
