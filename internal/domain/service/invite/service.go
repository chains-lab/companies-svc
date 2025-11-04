package invite

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
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
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (models.Employee, error)
	GetEmployeeByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (models.Employee, error)

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, ID uuid.UUID, answer string) error

	EmployeeExist(ctx context.Context, userID uuid.UUID) (bool, error)
}

type EventPublisher interface {
	PublishEmployeeCreated(ctx context.Context, employee models.Employee) error
	PublishInviteCreated(ctx context.Context, invite models.Invite) error
}

func (s Service) companyIsActive(ctx context.Context, companyID uuid.UUID) error {
	dis, err := s.db.GetCompanyByID(ctx, companyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by ID, cause: %w", err),
		)
	}
	if dis.IsNil() {
		return errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with ID %s not found", companyID),
		)
	}
	if dis.Status != enum.CompanyStatusActive {
		return errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("company with ID %s is not active", companyID),
		)
	}

	return nil
}
