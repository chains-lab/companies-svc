package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/distributors-svc/internal"
	"github.com/chains-lab/distributors-svc/internal/data"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/distributors-svc/internal/infra/jwtmanager"
	"github.com/chains-lab/distributors-svc/internal/rest"
	"github.com/chains-lab/distributors-svc/internal/rest/controller"
	"github.com/chains-lab/logium"
)

func Start(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(pg)

	jwtInviteManager := jwtmanager.NewManager(cfg)
	distributorSvc := distributor.NewService(database)
	employeeSvc := employee.NewService(database, jwtInviteManager)

	ctrl := controller.New(log, distributorSvc, employeeSvc)

	run(func() { rest.Run(ctx, cfg, log, ctrl) })
}
