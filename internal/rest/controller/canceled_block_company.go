package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) CanceledCompanyBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid company_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"company_id": err,
		})...)

		return
	}

	comp, err := a.domain.block.Cancel(r.Context(), companyID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to canceled block company")
		switch {
		case errors.Is(err, errx.ErrorcompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("company block %s canceled successfully by user %s", companyID, initiator.ID)

	ape.Render(w, http.StatusOK, responses.Company(comp))
}
