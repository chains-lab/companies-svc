package cmd

import (
	"context"
	"sync"

	"github.com/chains-lab/distributors-svc/internal"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/rest"
	"github.com/chains-lab/distributors-svc/internal/rest/handlers"
	"github.com/chains-lab/logium"
)

func Start(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	restSVC := rest.NewRest(cfg, log)

	run(func() {
		handl := handlers.NewAdapter(cfg, log, app)

		restSVC.Router(ctx, cfg, handl)
	})
}
