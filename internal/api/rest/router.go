package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CanceledDistributorBlock(w http.ResponseWriter, r *http.Request)
	CreateDistributorBlock(w http.ResponseWriter, r *http.Request)
	CreateDistributor(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
	GetActiveDistributorBlock(w http.ResponseWriter, r *http.Request)
	GetDistributor(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	ListBlockages(w http.ResponseWriter, r *http.Request)
	ListDistributors(w http.ResponseWriter, r *http.Request)
	ListEmployees(w http.ResponseWriter, r *http.Request)
	SentInvite(w http.ResponseWriter, r *http.Request)
	UpdateDistributor(w http.ResponseWriter, r *http.Request)
	UpdateDistributorStatus(w http.ResponseWriter, r *http.Request)
	AnswerToInvite(w http.ResponseWriter, r *http.Request)
	GetBlock(w http.ResponseWriter, r *http.Request)
}

func (a *Rest) Router(ctx context.Context, cfg config.Config, h Handlers) {
	svc := mdlv.ServiceGrant(enum.CitiesSVC, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin:     true,
		roles.SuperUser: true,
	})

	a.log.WithField("module", "api").Info("Starting API server")

	a.router.Route("/distributor-svc/", func(r chi.Router) {
		r.Use(svc)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/distributors", func(r chi.Router) {
				r.Get("/", h.ListDistributors)
				r.With(auth).Post("/", h.CreateDistributor)

				r.Route("/{distributor_id}", func(r chi.Router) {
					r.Get("/", h.GetDistributor)
					r.With(auth).Post("/", h.UpdateDistributor)

					r.Route("/status", func(r chi.Router) {
						r.With(auth).Post("/", h.UpdateDistributorStatus)
					})
				})
			})

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.ListEmployees)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetEmployee)
					r.With(auth).Delete("/", h.DeleteEmployee)
				})
			})

			r.With(auth).Route("/invite", func(r chi.Router) {
				r.Post("/", h.SentInvite)
				r.Post("/{token}", h.AnswerToInvite)
			})

			r.Route("/blocks", func(r chi.Router) {
				r.Get("/", h.ListBlockages)
				r.With(auth).With(sysadmin).Post("/", h.CreateDistributorBlock)
				r.Get("/active", h.GetActiveDistributorBlock)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", h.GetBlock)
					r.With(auth).With(sysadmin).Post("/", h.CanceledDistributorBlock)
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
