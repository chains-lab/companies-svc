package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetDistributor(w http.ResponseWriter, r *http.Request) {
	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor ID format")

		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
		return
	}

	distributor, err := s.app.GetDistributor(r.Context(), distributorID)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get distributor")

		switch {
		case errors.Is(err, errx.DistributorNotFound):
			ape.RenderErr(w, problems.NotFound("Distributor not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, distributor)
}
