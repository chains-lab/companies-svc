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

func (a Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	blockID, err := uuid.Parse(chi.URLParam(r, "block_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid block ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"block_id": err,
		})...)

		return
	}

	block, err := a.domain.company.GetBlock(r.Context(), blockID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to get block, ID: %s", blockID)
		switch {
		case errors.Is(err, errx.ErrorcompanyBlockNotFound):
			ape.RenderErr(w, problems.NotFound("CreteBlock not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CompanyBlock(block))
}
