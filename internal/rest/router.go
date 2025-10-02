package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/chains-lab/distributors-svc/internal"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
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
	CreateInvite(w http.ResponseWriter, r *http.Request)
	UpdateDistributor(w http.ResponseWriter, r *http.Request)
	UpdateDistributorStatus(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	GetBlock(w http.ResponseWriter, r *http.Request)
}

func (s *Service) Router(ctx context.Context, cfg internal.Config, h Handlers) {
	svc := mdlv.ServiceGrant(enum.CitiesSVC, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin:     true,
		roles.SuperUser: true,
	})
	user := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.User: true,
	})

	s.log.WithField("module", "api").Info("Starting API server")

	s.router.Route("/distributor-svc/", func(r chi.Router) {
		r.Use(svc)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/distributors", func(r chi.Router) {
				r.Get("/", h.ListDistributors)
				r.With(auth, user).Post("/", h.CreateDistributor)

				r.Route("/{distributor_id}", func(r chi.Router) {
					r.Get("/", h.GetDistributor)
					r.With(auth, user).Post("/", h.UpdateDistributor)

					r.Route("/status", func(r chi.Router) {
						r.With(auth, user).Post("/", h.UpdateDistributorStatus)
						r.With(auth, sysadmin).Post("/", h.CreateDistributorBlock)
					})
				})
			})

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.ListEmployees)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetEmployee)
					r.With(auth, user).Delete("/", h.DeleteEmployee)
				})
			})

			r.With(auth, user).Route("/invite", func(r chi.Router) {
				r.Post("/", h.CreateInvite)
				r.Post("/{token}", h.AcceptInvite)
			})

			r.Route("/blocks", func(r chi.Router) {
				r.Get("/", h.ListBlockages)
				r.With(auth, sysadmin).Post("/", h.CreateDistributorBlock)
				r.Get("/active", h.GetActiveDistributorBlock)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", h.GetBlock)
					r.With(auth, sysadmin).Post("/", h.CanceledDistributorBlock)
				})
			})
		})
	})

	s.Start(ctx)

	<-ctx.Done()
	s.Stop(ctx)
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		s.log.Infof("Starting server on port %s", s.cfg.Server.Port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (s *Service) Stop(ctx context.Context) {
	s.log.Info("Shutting down server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Errorf("Server shutdown failed: %v", err)
	}
}
