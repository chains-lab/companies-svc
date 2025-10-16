package employee

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UserGuesser interface {
	Guess(ctx context.Context, userIDs ...uuid.UUID) (map[uuid.UUID]models.Profile, error)
}

type Service struct {
	db          database
	userGuesser UserGuesser
}

func NewService(db database, userGuesser UserGuesser) Service {
	return Service{
		db:          db,
		userGuesser: userGuesser,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	CreateEmployee(ctx context.Context, input models.Employee) error
	FilterEmployees(ctx context.Context, filters FilterParams, page, size uint64) (models.EmployeeCollection, error)

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

	GetEmployee(
		ctx context.Context,
		params GetParams,
	) (models.Employee, error)

	UpdateEmployeeRole(ctx context.Context, userID uuid.UUID, newRole string, updatedAt time.Time) error
	DeleteEmployee(ctx context.Context, userID, companyID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, ID, UserID uuid.UUID, status string, acceptedAt time.Time) error
}
