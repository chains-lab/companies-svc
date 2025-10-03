package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) GetDistributor(w http.ResponseWriter, r *http.Request) {
	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"distributor_id": err,
		})...)

		return
	}

	distributor, err := a.domain.distributor.Get(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Error("failed to get distributor")
		switch {
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("DistributorID not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, distributor)
}
