package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetActiveDistributorBlock(w http.ResponseWriter, r *http.Request) {
	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor ID format")
		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))

		return
	}

	block, err := a.app.GetActiveDistributorBlock(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to get distributor %s active block", distributorID)
		switch {
		case errors.Is(err, errx.ErrorNoActiveBlockForDistributor):
			ape.RenderErr(w, problems.NotFound("no active block for distributor"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlock(block))
}
