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

func (a Adapter) CanceledDistributorBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid distributor_id")
		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))

		return
	}

	if initiator.Role != roles.Admin && initiator.Role != roles.SuperUser {
		a.log.Warnf("user %s with role %s attempted to canceled block %s", initiator.ID, initiator.Role, distributorID)
		ape.RenderErr(w, problems.Forbidden("only admin and superuser can unblock distributor"))

		return
	}

	dis, err := a.app.UnblockDistributor(r.Context(), distributorID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to canceled block distributor")
		switch {
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		case errors.Is(err, errx.ErrorNoActiveBlockForDistributor):
			ape.RenderErr(w, problems.Conflict("no active block for distributor"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor block %s canceled successfully by user %s", distributorID, initiator.ID)

	ape.Render(w, http.StatusOK, responses.Distributor(dis))
}
