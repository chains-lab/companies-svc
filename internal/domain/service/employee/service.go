package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type JwtManager interface {
	CreateInviteToken(
		inviteID uuid.UUID,
		role string,
		cityID uuid.UUID,
		ExpiredAt time.Time,
	) (string, error)

	DecryptInviteToken(tokenStr string) (models.InviteTokenData, error)

	HashInviteToken(tokenStr string) (string, error)
	VerifyInviteToken(tokenStr, hashed string) error
}

type Service struct {
	db  database
	jwt JwtManager
}

func NewService(db database, jwt JwtManager) Service {
	return Service{
		db:  db,
		jwt: jwt,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetDistributorByID(ctx context.Context, ID uuid.UUID) (models.Distributor, error)

	CreateEmployee(ctx context.Context, input models.Employee) error
	FilterEmployees(ctx context.Context, filters Filter, page, size uint64) (models.EmployeeCollection, error)
	GetEmployee(ctx context.Context, filters GetFilters) (models.Employee, error)
	UpdateEmployeeRole(ctx context.Context, userID uuid.UUID, newRole string, updatedAt time.Time) error
	DeleteEmployee(ctx context.Context, userID, distributorID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, ID, UserID uuid.UUID, status string, acceptedAt time.Time) error
}

func (s Service) DistributorIsActive(ctx context.Context, distributorID uuid.UUID) error {
	dis, err := s.db.GetDistributorByID(ctx, distributorID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get distributor by ID, cause: %w", err),
		)
	}
	if dis.IsNil() {
		return errx.ErrorDistributorNotFound.Raise(
			fmt.Errorf("distributor with ID %s not found", distributorID),
		)
	}
	if dis.Status != enum.DistributorStatusActive {
		return errx.ErrorDistributorIsNotActive.Raise(
			fmt.Errorf("distributor with ID %s is not active", distributorID),
		)
	}

	return nil
}
