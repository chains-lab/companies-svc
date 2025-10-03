package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) GetActiveDistributorBlock(w http.ResponseWriter, r *http.Request) {
	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"distributor_id": err,
		})...)

		return
	}

	block, err := a.domain.distributor.GetActiveDistributorBlock(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to get distributor %s active block", distributorID)
		switch {
		case errors.Is(err, errx.ErrorDistributorBlockNotFound):
			ape.RenderErr(w, problems.NotFound("no active block for distributor"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlock(block))
}
