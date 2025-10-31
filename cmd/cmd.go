package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/companies-svc/internal"
	"github.com/chains-lab/companies-svc/internal/data"
	"github.com/chains-lab/companies-svc/internal/domain/service/block"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/domain/service/invite"
	"github.com/chains-lab/companies-svc/internal/events/publisher"
	"github.com/chains-lab/companies-svc/internal/jwtmanager"
	"github.com/chains-lab/companies-svc/internal/rest"
	"github.com/chains-lab/companies-svc/internal/rest/controller"
	"github.com/chains-lab/companies-svc/internal/rest/middlewares"
	"github.com/chains-lab/companies-svc/internal/usrguesser"
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
	userGuesser := usrguesser.NewService(cfg.Profile.Url, nil)

	eventWriter := publisher.New(cfg.Kafka.Broker)

	companiesSvc := company.NewService(database, eventWriter)
	employeeSvc := employee.NewService(database, userGuesser, eventWriter)
	inviteSvc := invite.NewService(database, jwtInviteManager, eventWriter)
	blockSvc := block.NewService(database)

	ctrl := controller.New(log, companiesSvc, employeeSvc, inviteSvc, blockSvc)
	mdlv := middlewares.New(log)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
