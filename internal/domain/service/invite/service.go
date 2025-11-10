package invite

import (
	"context"

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
	GetCompanyEmployees(ctx context.Context, companyID uuid.UUID, roles ...string) (models.EmployeesCollection, error)

	EmployeeExist(ctx context.Context, userID uuid.UUID) (bool, error)
}

type EventPublisher interface {
	PublishInviteAccepted(
		ctx context.Context,
		invite models.Invite,
		company models.Company,
		recipients []uuid.UUID,
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
		recipients []uuid.UUID,
	) error

	PublishEmployeeCreated(
		ctx context.Context,
		company models.Company,
		employee models.Employee,
		recipients []uuid.UUID,
	) error
}
