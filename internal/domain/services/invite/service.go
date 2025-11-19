package invite

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db    database
	event EventPublisher
}

func NewService(db database, publisher EventPublisher) Service {
	return Service{
		db:    db,
		event: publisher,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	CreateEmployee(ctx context.Context, input models.Employee) error

	GetEmployee(ctx context.Context, userID, companyID uuid.UUID) (models.Employee, error)
	GetCompanyOwner(ctx context.Context, companyID uuid.UUID) (models.Employee, error)

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, ID uuid.UUID, reply string) error
	GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error)
}

type EventPublisher interface {
	PublishInviteAccepted(
		ctx context.Context,
		invite models.Invite,
		company models.Company,
		recipients ...uuid.UUID,
	) error

	PublishInviteDeclined(
		ctx context.Context,
		invite models.Invite,
		company models.Company,
	) error

	PublishInviteCreated(
		ctx context.Context,
		invite models.Invite,
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

func (s Service) getCompany(
	ctx context.Context,
	companyID uuid.UUID,
) (models.Company, error) {
	company, err := s.db.GetCompanyByID(ctx, companyID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company with id %s: %w", companyID, err),
		)
	}
	if company.IsNil() {
		return models.Company{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with id %s not found", companyID),
		)
	}

	return company, nil
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
