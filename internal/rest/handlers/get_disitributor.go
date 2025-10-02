package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetDistributor(w http.ResponseWriter, r *http.Request) {
	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor ID format")
		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))

		return
	}

	distributor, err := a.app.GetDistributor(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Error("failed to get distributor")
		switch {
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("Distributor not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, distributor)
}
