package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/companies-svc/internal"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CanceledCompanyBlock(w http.ResponseWriter, r *http.Request)
	CreateCompanyBlock(w http.ResponseWriter, r *http.Request)
	CreateCompany(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
	GetActiveCompanyBlock(w http.ResponseWriter, r *http.Request)
	GetCompany(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	FilterBlockages(w http.ResponseWriter, r *http.Request)
	FilterCompanies(w http.ResponseWriter, r *http.Request)
	ListEmployees(w http.ResponseWriter, r *http.Request)
	CreateInvite(w http.ResponseWriter, r *http.Request)
	UpdateCompany(w http.ResponseWriter, r *http.Request)
	UpdateCompaniesStatus(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
	GetBlock(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	ServiceGrant(serviceName, skService string) func(http.Handler) http.Handler
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middlewares, h Handlers) {
	svc := m.ServiceGrant(cfg.Service.Name, cfg.JWT.Service.SecretKey)
	auth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})
	user := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.User: true,
	})

	r := chi.NewRouter()

	log.WithField("module", "api").Info("Starting API server")

	r.Route("/company-svc/", func(r chi.Router) {
		r.Use(svc)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/companies", func(r chi.Router) {
				r.Get("/", h.FilterCompanies)
				r.With(auth, user).Post("/", h.CreateCompany)

				r.Route("/{company_id}", func(r chi.Router) {
					r.Get("/", h.GetCompany)
					r.With(auth, user).Post("/", h.UpdateCompany)

					r.Route("/status", func(r chi.Router) {
						r.With(auth, user).Post("/", h.UpdateCompaniesStatus)
						r.With(auth, sysadmin).Post("/", h.CreateCompanyBlock)
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
				r.Get("/", h.FilterBlockages)
				r.With(auth, sysadmin).Post("/", h.CreateCompanyBlock)
				r.Get("/active", h.GetActiveCompanyBlock)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", h.GetBlock)
					r.With(auth, sysadmin).Post("/", h.CanceledCompanyBlock)
				})
			})
		})
	})

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	<-ctx.Done()

	log.Info("shutting down REST service")
}
