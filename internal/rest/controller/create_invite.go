package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
)

func (a Service) CreateInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.CreateEmployeeInvite(r)
	if err != nil {
		a.log.WithError(err).Errorf("invalid create employee invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	invite, err := a.domain.employee.CreateInvite(r.Context(), initiator.ID, employee.SentInviteParams{
		DistributorID: req.Data.Attributes.DistributorId,
		Role:          req.Data.Attributes.Role,
	})
	if err != nil {
		a.log.WithError(err).Errorf("failed to create employee invite")
		switch {
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		case errors.Is(err, errx.ErrorDistributorIsBlocked):
			ape.RenderErr(w, problems.Conflict("distributor is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invites(invite))
}
