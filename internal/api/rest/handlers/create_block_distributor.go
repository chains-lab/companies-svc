package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"

	"github.com/chains-lab/distributors-svc/internal/api/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
)

func (s Service) CreateDistributorBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	req, err := requests.CreateDistributorBlock(r)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to decode block distributor request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	distributorID, err := uuid.Parse(req.Data.Attributes.DistributorId)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor id: %s", req.Data.Attributes.DistributorId)

		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
		return
	}

	block, err := s.app.BlockDistributor(r.Context(), initiator.ID, distributorID, req.Data.Attributes.Reason)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to block distributor")

		switch {
		case errors.Is(err, errx.DistributorHaveAlreadyActiveBlock):
			ape.RenderErr(w, problems.Conflict("distributor already have active block"))
		case errors.Is(err, errx.DistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("distributor %s blocked successfully by user %s", distributorID, initiator.ID)

	responses.DistributorBlock(block)
}
