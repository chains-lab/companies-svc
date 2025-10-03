package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/companies-svc/internal/data/pgdb"
)

type Database struct {
	sql SqlDB
}

func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.companies.New().Transaction(ctx, fn)
}

type SqlDB struct {
	companies pgdb.CompaniesQ
	blockages pgdb.BlockagesQ
	employees pgdb.EmployeesQ
	invites   pgdb.InvitesQ
}

func NewDatabase(db *sql.DB) *Database {
	companiesql := pgdb.NewcompaniesQ(db)
	blockagesSql := pgdb.NewBlockagesQ(db)
	employeesSql := pgdb.NewEmployeesQ(db)
	invitesSql := pgdb.NewInvitesQ(db)

	return &Database{
		sql: SqlDB{
			companies: companiesql,
			blockages: blockagesSql,
			employees: employeesSql,
			invites:   invitesSql,
		},
	}
}
