package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
)

func (a Service) UpdateDistributor(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateDistributor(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update distributor request")

		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	input := distributor.UpdateParams{}
	if req.Data.Attributes.Name != nil {
		input.Name = req.Data.Attributes.Name
	}
	if req.Data.Attributes.Icon != nil {
		input.Icon = req.Data.Attributes.Icon
	}

	res, err := a.domain.distributor.Update(r.Context(), req.Data.Id, input)
	if err != nil {
		a.log.WithError(err).Errorf("failed to update distributor name for ID: %s", req.Data.Id)
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

	a.log.Infof("distributor %s updated successfully", fmt.Sprint(res.ID))

	ape.Render(w, http.StatusOK, responses.Distributor(res))
}
