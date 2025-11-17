package employee

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

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	CreateEmployee(ctx context.Context, input models.Employee) error

	GetEmployee(ctx context.Context, ID uuid.UUID) (models.Employee, error)
	GetEmployeeUserInCompany(ctx context.Context, userID, companyID uuid.UUID) (models.Employee, error)
	GetCompanyOwner(ctx context.Context, companyID uuid.UUID) (models.Employee, error)

	FilterEmployees(ctx context.Context, filters FilterParams, page, size uint64) (models.EmployeesCollection, error)

	GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error)

	UpdateEmployee(ctx context.Context, userID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	DeleteEmployee(ctx context.Context, ID uuid.UUID) error
}

type EventPublisher interface {
	PublishEmployeeCreated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients ...uuid.UUID,
	) error

	PublishEmployeeUpdated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients ...uuid.UUID,
	) error

	PublishEmployeeDeleted(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		Recipients ...uuid.UUID,
	) error
}

func (s Service) validateInitiator(
	ctx context.Context,
	userID uuid.UUID,
	companyID uuid.UUID,
	roles ...string,
) (models.Employee, error) {
	employee, err := s.db.GetEmployeeUserInCompany(ctx, userID, companyID)
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

func (s Service) getCompany(ctx context.Context, companyID uuid.UUID) (models.Company, error) {
	company, err := s.db.GetCompanyByID(ctx, companyID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by EmployeeID %s, cause: %w", companyID, err),
		)
	}
	if company.IsNil() {
		return models.Company{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with EmployeeID %s not found", companyID),
		)
	}

	return company, nil
}
