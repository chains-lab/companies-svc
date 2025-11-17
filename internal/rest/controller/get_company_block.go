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

func (s Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	blockID, err := uuid.Parse(chi.URLParam(r, "block_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid block EmployeeID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"block_id": err,
		})...)

		return
	}

	block, err := s.domain.block.Get(r.Context(), blockID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get block, EmployeeID: %s", blockID)
		switch {
		case errors.Is(err, errx.ErrorCompanyBlockNotFound):
			ape.RenderErr(w, problems.NotFound("Crete not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CompanyBlock(block))
}
