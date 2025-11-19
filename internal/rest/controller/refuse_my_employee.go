package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) RefuseMyEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid company EmployeeID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"company_id": err,
		})...)

		return
	}

	err = s.domain.employee.DeleteMe(r.Context(), initiator.ID, companyID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get employee")
		switch {
		case errors.Is(err, errx.ErrorOwnerCannotRefuseSelf):
			ape.RenderErr(w, problems.Forbidden("owner cannot refuse himself"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
