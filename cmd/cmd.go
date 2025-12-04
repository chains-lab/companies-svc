package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/organizations-svc/internal"
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

	database := repo.NewDatabase(pg)

	eventPublisher := publisher.New(cfg.Kafka.Broker)

	organizationsSvc := company.NewService(database, eventPublisher)
	employeeSvc := employee.NewService(database, eventPublisher)
	inviteSvc := invite.NewService(database, eventPublisher)
	blockSvc := block.NewService(database, eventPublisher)

	ctrl := controller.New(log, organizationsSvc, employeeSvc, inviteSvc, blockSvc)
	mdlv := middlewares.New(log)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
