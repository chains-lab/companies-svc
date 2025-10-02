package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdateDistributor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")

		return
	}

	req, err := requests.UpdateDistributor(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update distributor request")

		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	if req.Data.Id != chi.URLParam(r, "distributor_id") {
		ape.RenderErr(w,
			problems.InvalidParameter("distributor_id", fmt.Errorf("path ID and body ID do not match")),
			problems.InvalidPointer("/data/id", fmt.Errorf("path ID and body ID do not match")),
		)

		return
	}

	distributorID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor id: %s", req.Data.Id)
		ape.RenderErr(w,
			problems.InvalidParameter("distributor_id", err),
			problems.InvalidPointer("/data/id", err),
		)

		return
	}

	input := app.UpdateDistributorParams{}
	if req.Data.Attributes.Name != nil {
		input.Name = req.Data.Attributes.Name
	}
	if req.Data.Attributes.Icon != nil {
		input.Icon = req.Data.Attributes.Icon
	}

	distributor, err := a.app.UpdateDistributor(r.Context(), initiator.ID, distributorID, input)
	if err != nil {
		a.log.WithError(err).Errorf("failed to update distributor name for ID: %s", distributorID)
		switch {
		case errors.Is(err, errx.ErrorInitiatorNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorEmployeeHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee has not enough rights"))
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		case errors.Is(err, errx.ErrorDistributorIsBlocked):
			ape.RenderErr(w, problems.Conflict("distributor is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor %s updated successfully", distributorID)

	ape.Render(w, http.StatusOK, responses.Distributor(distributor))
}
