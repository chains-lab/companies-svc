package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
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

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)

	CreateEmployee(ctx context.Context, input models.Employee) error

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

	UpdateEmployeeRole(ctx context.Context, userID uuid.UUID, newRole string, updatedAt time.Time) error
	DeleteEmployee(ctx context.Context, userID, companyID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, ID, UserID uuid.UUID, status string, acceptedAt time.Time) error
}

func (s Service) companyIsActive(ctx context.Context, companyID uuid.UUID) error {
	dis, err := s.db.GetCompanyByID(ctx, companyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by ID, cause: %w", err),
		)
	}
	if dis.IsNil() {
		return errx.ErrorcompanyNotFound.Raise(
			fmt.Errorf("company with ID %s not found", companyID),
		)
	}
	if dis.Status != enum.DistributorStatusActive {
		return errx.ErrorcompanyIsNotActive.Raise(
			fmt.Errorf("company with ID %s is not active", companyID),
		)
	}

	return nil
}
