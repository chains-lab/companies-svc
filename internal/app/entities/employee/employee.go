package employee

import (
	"context"
	"database/sql"

	"github.com/chains-lab/distributors-svc/internal/app/jwtmanager"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/google/uuid"
)

type employeesQ interface {
	New() dbx.EmployeeQ

	Insert(ctx context.Context, input dbx.Employee) error
	Get(ctx context.Context) (dbx.Employee, error)
	Select(ctx context.Context) ([]dbx.Employee, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterUserID(id uuid.UUID) dbx.EmployeeQ
	FilterDistributorID(id ...uuid.UUID) dbx.EmployeeQ
	FilterRole(role ...string) dbx.EmployeeQ

	OrderByRole(ascend bool) dbx.EmployeeQ

	Page(limit, offset uint64) dbx.EmployeeQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

type inviteQ interface {
	New() dbx.InviteQ

	Insert(ctx context.Context, input dbx.Invite) error
	Get(ctx context.Context) (dbx.Invite, error)
	Select(ctx context.Context) ([]dbx.Invite, error)
	Update(ctx context.Context, params dbx.UpdateInviteParams) error
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.InviteQ
	FilterDistributorID(distributorID uuid.UUID) dbx.InviteQ
	FilterUserID(userID uuid.UUID) dbx.InviteQ
	FilterRole(role ...string) dbx.InviteQ
	FilterStatus(status ...string) dbx.InviteQ

	OrderByCreatedAt(asc bool) dbx.InviteQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.InviteQ
}

type Employee struct {
	employee employeesQ
	invite   inviteQ
	jwt      jwtmanager.Manager
}

func NewEmployee(db *sql.DB, cfg config.Config) Employee {
	return Employee{
		employee: dbx.NewEmployeesQ(db),
		invite:   dbx.NewInviteQ(db),
		jwt:      jwtmanager.NewManager(cfg),
	}
}

func inviteFromDB(inv dbx.Invite, token string) models.Invite {
	res := models.Invite{
		ID:            inv.ID,
		Status:        inv.Status,
		Role:          inv.Role,
		DistributorID: inv.DistributorID,
		Token:         token,
		ExpiresAt:     inv.ExpiresAt,
		CreatedAt:     inv.CreatedAt,
	}
	if inv.UserID.Valid {
		res.UserID = &inv.UserID.UUID
	}
	if inv.AnsweredAt.Valid {
		res.AnsweredAt = &inv.AnsweredAt.Time
	}

	return res
}
