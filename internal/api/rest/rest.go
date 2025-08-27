package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/handlers"
	"github.com/chains-lab/distributors-svc/internal/api/rest/mdlv"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, logger, app)

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Rest) Run(ctx context.Context) {
	userAuth := mdlv.AuthMdl(a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(roles.Admin, roles.SuperUser)
	svcAuth := mdlv.ServiceAuthMdl(a.cfg.JWT.Service.SecretKey)

	a.log.WithField("module", "api").Info("Starting API server")

	a.router.Route("/distributor-svc/", func(r chi.Router) {
		r.Use(svcAuth)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/distributors", func(r chi.Router) {
				r.Get("/", a.handlers.SelectDistributors)
				r.With(userAuth).Post("/", a.handlers.CreateDistributor)

				r.Route("/{distributor_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetDistributor)
					r.With(userAuth).Post("/", a.handlers.UpdateDistributor)

					r.Route("/status", func(r chi.Router) {
						r.With(userAuth).Post("/", a.handlers.UpdateDistributorStatus)
					})
				})
			})

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", a.handlers.SelectEmployees)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetEmployee)
					r.With(userAuth).Delete("/", a.handlers.DeleteEmployee)

					r.With(userAuth).Put("/role", a.handlers.UpdateEmployeeRole)
				})
			})

			r.Route("/invites", func(r chi.Router) {
				r.Get("/", a.handlers.SelectInvites)
				r.With(userAuth).Post("/", a.handlers.CreateInvite)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetInvite)

					r.With(userAuth).Post("/{status}", a.handlers.InteractToInvite)
				})

			})

			r.Route("/blocks", func(r chi.Router) {
				r.Get("/", a.handlers.SelectDistributorBlocks)
				r.With(userAuth).With(adminGrant).Post("/", a.handlers.CreateDistributorBlock)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetDistributorBlock)
					r.With(userAuth).With(adminGrant).Post("/", a.handlers.CanceledDistributorBlock)
				})
			})
		})
	})

	a.Start(ctx)

	<-ctx.Done()
	a.Stop(ctx)
}

func (a *Rest) Start(ctx context.Context) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Errorf("Server shutdown failed: %v", err)
	}
}
