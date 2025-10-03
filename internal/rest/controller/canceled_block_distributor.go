package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) CanceledDistributorBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid distributor_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"distributor_id": err,
		})...)

		return
	}

	dis, err := a.domain.distributor.CancelBlock(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to canceled block distributor")
		switch {
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor block %s canceled successfully by user %s", distributorID, initiator.ID)

	ape.Render(w, http.StatusOK, responses.Distributor(dis))
}
