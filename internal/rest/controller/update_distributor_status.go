package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
)

func (a Service) UpdateDistributorStatus(w http.ResponseWriter, r *http.Request) {
	//initiator, err := meta.User(r.Context())
	//if err != nil {
	//	a.log.WithError(err).Error("failed to get user from context")
	//	ape.RenderErr(w, problems.Unauthorized("user not found in context"))
	//	return
	//}

	req, err := requests.UpdateDistributorStatus(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update distributor status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := a.domain.distributor.UpdateStatus(r.Context(), req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		a.log.WithError(err).Errorf("failed to set distributor %s status to active", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
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

	a.log.Infof("distributor %s status set to active successfully", res.ID)

	ape.Render(w, http.StatusOK, responses.Distributor(res))
}
