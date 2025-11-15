package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) GetActiveCompanyBlock(w http.ResponseWriter, r *http.Request) {
	companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid company ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"company_id": err,
		})...)

		return
	}

	block, err := s.domain.block.GetActiveCompanyBlock(r.Context(), companyID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get company %s active block", companyID)
		switch {
		case errors.Is(err, errx.ErrorCompanyBlockNotFound):
			ape.RenderErr(w, problems.NotFound("no active block for company"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CompanyBlock(block))
}
