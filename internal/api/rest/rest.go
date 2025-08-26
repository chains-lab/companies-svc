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
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log *logrus.Entry
	cfg config.Config
}

func NewRest(cfg config.Config, log *logrus.Logger, app app.App) Rest {
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

func (a *Rest) Run(ctx context.Context, log *logrus.Logger) {
	userAuth := mdlv.AuthMdl(a.cfg.JWT.User.AccessToken.SecretKey)
	adminGrant := mdlv.AccessGrant(roles.Admin, roles.SuperUser)
	svcAuth := mdlv.ServiceAuthMdl(a.cfg.JWT.Service.SecretKey)
	requestID := mdlv.RequestIDMdl()

	a.log.WithField("module", "api").Info("Starting API server")

	a.router.Route("/distributor-svc/", func(r chi.Router) {
		r.Use(requestID)
		r.Use(svcAuth)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/distributors", func(r chi.Router) {
				r.Route("/{distributor_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetDistributor)
					r.Post("/update", a.handlers.UpdateDistributor)

					r.Route("/status", func(r chi.Router) {
						r.Use(userAuth)

						r.Post("/activate", a.handlers.SetDistributorStatusActive)
						r.Post("/deactivate", a.handlers.SetDistributorStatusInactive)
						r.With(adminGrant).Post("/block", a.handlers.BlockDistributor)
					})
				})
			})

			r.Route("employees", func(r chi.Router) {
				r.Get("/", a.handlers.SelectEmployees)

				r.Route("/{employee_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetEmployee)

					r.With(userAuth).Put("/", a.handlers.UpdateEmployee)
					r.With(userAuth).Delete("/", a.handlers.DeleteEmployee)
				})
			})

			r.Route("/invites", func(r chi.Router) {
				r.With(userAuth).Post("/", a.handlers.SendInvite)

				r.Route("/{invite_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetInvite)

					r.With(userAuth).Post("/accept", a.handlers.AcceptInvite)
					r.With(userAuth).Post("/reject", a.handlers.RejectInvite)
					r.With(userAuth).Post("/withdrew", a.handlers.WithdrewInvite)
				})

			})

			r.Route("/blocks", func(r chi.Router) {
				r.Get("/", a.handlers.SelectDistributorBlocks)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetDistributorBlock)

					r.Post("/unlock", a.handlers.UnblockDistributor)
				})
			})
		})
	})

	a.Start(ctx, log)

	<-ctx.Done()
	a.Stop(ctx, log)
}

func (a *Rest) Start(ctx context.Context, log *logrus.Logger) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Rest) Stop(ctx context.Context, log *logrus.Logger) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	}
}
