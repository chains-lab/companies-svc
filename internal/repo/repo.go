package repo

import (
	"context"
	"database/sql"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/repo/pgdb"
)

type Repo struct {
	sql SqlDB
}

func (r *Repo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.sql.companies.New().Transaction(ctx, fn)
}

type SqlDB struct {
	companies pgdb.CompaniesQ
	blockages pgdb.BlocksQ
	employees pgdb.EmployeesQ
	invites   pgdb.InvitesQ
}

func NewDatabase(db *sql.DB) *Repo {
	return &Repo{
		sql: SqlDB{
			companies: pgdb.NewcompaniesQ(db),
			blockages: pgdb.NewBlocksQ(db),
			employees: pgdb.NewEmployeesQ(db),
			invites:   pgdb.NewInvitesQ(db),
		},
	}
}

func companyModelToSchema(model models.Company) pgdb.Company {
	return pgdb.Company{
		ID:        model.ID,
		Name:      model.Name,
		Icon:      model.Icon,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func companiesSchemaToModel(schema pgdb.Company) models.Company {
	return models.Company{
		ID:        schema.ID,
		Name:      schema.Name,
		Icon:      schema.Icon,
		Status:    schema.Status,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
}

func blockModelToSchema(m models.CompanyBlock) pgdb.CompanyBlock {
	block := pgdb.CompanyBlock{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		InitiatorID: m.InitiatorID,
		Reason:      m.Reason,
		Status:      m.Status,
		BlockedAt:   m.BlockedAt,
	}
	if m.CanceledAt != nil {
		block.CanceledAt = m.CanceledAt
	}

	return block
}

func companyBlockSchemaToModel(s pgdb.CompanyBlock) models.CompanyBlock {
	block := models.CompanyBlock{
		ID:          s.ID,
		CompanyID:   s.CompanyID,
		InitiatorID: s.InitiatorID,
		Reason:      s.Reason,
		Status:      s.Status,
		BlockedAt:   s.BlockedAt,
	}
	if s.CanceledAt != nil {
		block.CanceledAt = s.CanceledAt
	}

	return block
}

func inviteModelToSchema(m models.Invite) pgdb.Invite {
	res := pgdb.Invite{
		ID:        m.ID,
		Status:    m.Status,
		Role:      m.Role,
		CompanyID: m.CompanyID,
		CreatedAt: m.CreatedAt,
		ExpiresAt: m.ExpiresAt,
	}

	return res
}

func inviteSchemaToModel(m pgdb.Invite) models.Invite {
	res := models.Invite{
		ID:        m.ID,
		Status:    m.Status,
		Role:      m.Role,
		CompanyID: m.CompanyID,
		CreatedAt: m.CreatedAt,
		ExpiresAt: m.ExpiresAt,
	}

	return res
}

func employeeModelToSchema(input models.Employee) pgdb.Employee {
	return pgdb.Employee{
		UserID:    input.UserID,
		CompanyID: input.CompanyID,
		Role:      input.Role,
		UpdatedAt: input.UpdatedAt,
		CreatedAt: input.CreatedAt,
	}
}

func employeeSchemaToModel(input pgdb.Employee) models.Employee {
	return models.Employee{
		UserID:    input.UserID,
		CompanyID: input.CompanyID,
		Role:      input.Role,
		UpdatedAt: input.UpdatedAt,
		CreatedAt: input.CreatedAt,
	}
}
