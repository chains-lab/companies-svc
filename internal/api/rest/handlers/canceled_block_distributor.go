package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/gatekit/roles"
)

func (s Service) CanceledDistributorBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	blockID, err := uuid.Parse(chi.URLParam(r, "block_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor id: %s", chi.URLParam(r, "block_id"))

		ape.RenderErr(w, problems.InvalidParameter("block_id", err))
		return
	}

	if initiator.Role != roles.Admin && initiator.Role != roles.SuperUser {
		s.Log(r).Warnf("user %s with role %s attempted to canceled block %s", initiator.ID, initiator.Role, blockID)

		ape.RenderErr(w, problems.Forbidden("only admin and superuser can unblock distributor"))
		return
	}

	block, err := s.app.UnblockDistributor(r.Context(), blockID)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to canceled block distributor")

		switch {
		case errors.Is(err, errx.DistributorBlockNotFound):
			ape.RenderErr(w, problems.NotFound("distributor block not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("distributor block %s canceled successfully by user %s", blockID, initiator.ID)

	ape.Render(w, http.StatusOK, responses.DistributorBlock(block))
}
