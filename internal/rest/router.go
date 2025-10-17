package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chains-lab/companies-svc/internal"
	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/restkit/roles"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	CreateCompanyBlock(w http.ResponseWriter, r *http.Request)
	GetBlock(w http.ResponseWriter, r *http.Request)
	GetActiveCompanyBlock(w http.ResponseWriter, r *http.Request)
	FilterBlockages(w http.ResponseWriter, r *http.Request)
	CanceledCompanyBlock(w http.ResponseWriter, r *http.Request)

	CreateCompany(w http.ResponseWriter, r *http.Request)
	GetCompany(w http.ResponseWriter, r *http.Request)
	FilterCompanies(w http.ResponseWriter, r *http.Request)
	UpdateCompany(w http.ResponseWriter, r *http.Request)
	UpdateCompaniesStatus(w http.ResponseWriter, r *http.Request)

	ListEmployees(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)

	GetMyEmployee(w http.ResponseWriter, r *http.Request)
	RefuseMyEmployee(w http.ResponseWriter, r *http.Request)

	CreateInvite(w http.ResponseWriter, r *http.Request)
	AcceptInvite(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler

	CompanyMember(
		UserCtxKey interface{},
		allowedCompanyRoles map[string]bool,
	) func(http.Handler) http.Handler

	CompanyMemberOrAdmin(
		UserCtxKey interface{},
		allowedCompanyRoles map[string]bool,
		allowedAdminRoles map[string]bool,
	) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middlewares, h Handlers) {
	auth := m.Auth(meta.UserCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysadmin := m.RoleGrant(meta.UserCtxKey, map[string]bool{
		roles.Admin: true,
	})

	companyAdmin := m.CompanyMember(meta.UserCtxKey, map[string]bool{
		enum.EmployeeRoleOwner: true,
		enum.EmployeeRoleAdmin: true,
	})
	companyMember := m.CompanyMember(meta.UserCtxKey, map[string]bool{
		enum.EmployeeRoleOwner:     true,
		enum.EmployeeRoleAdmin:     true,
		enum.EmployeeRoleModerator: true,
	})

	r := chi.NewRouter()

	log.WithField("module", "api").Info("Starting API server")

	r.Route("/company-svc/", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Route("/companies", func(r chi.Router) {
				r.Get("/", h.FilterCompanies)
				r.With(auth).Post("/", h.CreateCompany)

				r.Route("/{company_id}", func(r chi.Router) {
					r.Get("/", h.GetCompany)
					r.With(auth, companyAdmin).Post("/", h.UpdateCompany)

					r.With(auth, companyAdmin).Post("/status", h.UpdateCompaniesStatus)
					r.With(auth, sysadmin).Post("/block", h.CreateCompanyBlock)
				})
			})

			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.ListEmployees)

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetEmployee)
					r.With(auth, companyAdmin).Delete("/", h.DeleteEmployee)
				})

				r.With(auth, companyMember).Route("/me", func(r chi.Router) {
					r.Get("/", h.GetMyEmployee)
					r.Delete("/", h.RefuseMyEmployee)
				})
			})

			r.With(auth).Route("/invite", func(r chi.Router) {
				r.Post("/", h.CreateInvite)
				r.Post("/{token}", h.AcceptInvite)
			})

			r.Route("/blocks", func(r chi.Router) {
				r.Get("/", h.FilterBlockages)
				r.With(auth, sysadmin).Post("/", h.CreateCompanyBlock)
				r.Get("/{company_id}", h.GetActiveCompanyBlock)

				r.Route("/{block_id}", func(r chi.Router) {
					r.Get("/", h.GetBlock)
					r.With(auth, sysadmin).Post("/", h.CanceledCompanyBlock)
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Rest.Port,
		Handler:           r,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("REST shutdown error: %v", err)
	} else {
		log.Info("REST server stopped")
	}
}
