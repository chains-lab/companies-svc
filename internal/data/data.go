package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
)

type Database struct {
	sql SqlDB
}

func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.distributors.New().Transaction(ctx, fn)
}

type SqlDB struct {
	distributors pgdb.DistributorsQ
	blockages    pgdb.BlockagesQ
	employees    pgdb.EmployeesQ
	invites      pgdb.InvitesQ
}

func NewDatabase(db *sql.DB) *Database {
	distributorSql := pgdb.NewDistributorsQ(db)
	blockagesSql := pgdb.NewBlockagesQ(db)
	employeesSql := pgdb.NewEmployeesQ(db)
	invitesSql := pgdb.NewInvitesQ(db)

	return &Database{
		sql: SqlDB{
			distributors: distributorSql,
			blockages:    blockagesSql,
			employees:    employeesSql,
			invites:      invitesSql,
		},
	}
}
