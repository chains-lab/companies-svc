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
		ownerID uuid.UUID,
	) error

	PublishCompanyDeleted(
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

func (s Service) getEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	empl, err := s.db.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user ID, cause: %w", err),
		)
	}
	if empl.IsNil() {
		return models.Employee{}, errx.ErrorEmployeeNotFound.Raise(
			fmt.Errorf("employee for user ID %s not found", userID),
		)
	}

	return empl, nil
}

func (s Service) validateInitiatorRight(
	ctx context.Context,
	initiatorID uuid.UUID,
	companyID uuid.UUID,
	roles ...string,
) (models.Employee, error) {
	employee, err := s.db.GetEmployeeByUserID(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user ID, cause: %w", err),
		)
	}
	if employee.IsNil() {
		return models.Employee{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("employee for user ID %s not found", initiatorID),
		)
	}

	if employee.CompanyID != companyID {
		return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
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
			return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
				fmt.Errorf("initiator have not enough rights in company: %s", companyID),
			)
		}
	}

	return employee, nil
}
