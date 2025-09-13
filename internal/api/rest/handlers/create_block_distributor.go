package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/distributors-svc/internal/api/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
)

func (a Adapter) CreateDistributorBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateDistributorBlock(r)
	if err != nil {
		a.log.WithError(err).Error("failed to decode block distributor request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Attributes.DistributorId != chi.URLParam(r, "distributor_id") {
		ape.RenderErr(w,
			problems.InvalidParameter("distributor_id", fmt.Errorf("path ID and body ID do not match")),
			problems.InvalidPointer("/data/attributes.distributor_id", fmt.Errorf("path ID and body ID do not match")),
		)

		return
	}

	distributorID, err := uuid.Parse(req.Data.Attributes.DistributorId)
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor id: %s", req.Data.Attributes.DistributorId)
		ape.RenderErr(w,
			problems.InvalidParameter("distributor_id", err),
			problems.InvalidPointer("/data/attributes.distributor_id", err),
		)

		return
	}

	block, err := a.app.BlockDistributor(r.Context(), initiator.ID, distributorID, req.Data.Attributes.Reason)
	if err != nil {
		a.log.WithError(err).Errorf("failed to block distributor")
		switch {
		case errors.Is(err, errx.DistributorHaveAlreadyActiveBlock):
			ape.RenderErr(w, problems.Conflict("distributor already have active block"))
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor %s blocked successfully by user %s", distributorID, initiator.ID)

	responses.DistributorBlock(block)
}
