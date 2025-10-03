package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/distributors-svc/internal"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/logium"
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

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, h Handlers) {
	svc := mdlv.ServiceGrant(enum.CitiesSVC, cfg.JWT.Service.SecretKey)
	auth := mdlv.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})
	user := mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.User: true,
	})

	r := chi.NewRouter()

	log.WithField("module", "api").Info("Starting API server")

	r.Route("/distributor-svc/", func(r chi.Router) {
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

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
