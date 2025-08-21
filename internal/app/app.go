package app

import (
	"database/sql"

	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/distributors-svc/internal/dbx"
)

type App struct {
	distributor distributorsQ
	employee    employeesQ
	block       blockagesQ
	invite      inviteQ
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		distributor: dbx.NewDistributorsQ(pg),
		employee:    dbx.NewEmployeesQ(pg),
		block:       dbx.NewBlockagesQ(pg),
		invite:      dbx.NewInvitesQ(pg),
	}, nil
}
