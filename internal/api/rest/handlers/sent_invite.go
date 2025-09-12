package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
)

func (s Service) SentInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req, err := requests.CreateEmployeeInvite(r)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid create employee invite request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userID, err := uuid.Parse(req.Data.Attributes.UserId)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid user ID: %s", req.Data.Attributes.UserId)

		ape.RenderErr(w, problems.InvalidParameter("user_id", err))
		return
	}

	distributorID, err := uuid.Parse(req.Data.Attributes.DistributorId)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor ID: %s", req.Data.Attributes.DistributorId)

		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
		return
	}

	invite, err := s.app.SentInvite(r.Context(), app.SentInviteParams{
		InitiatorID: initiator.ID,
		Role:        req.Data.Attributes.Role,
	})
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to create employee invite")

		switch {
		case errors.Is(err, errx.ErrorInvalidEmployeeRole):
			ape.RenderErr(w, problems.InvalidParameter("role", err))
		case errors.Is(err, errx.ErrorInitiatorRoleHaveNotEnoughRights):
			ape.RenderErr(w, problems.PreconditionFailed("initiator have not enough rights to invite this role"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("invite sent successfully for user ID: %s in distributor ID: %s", userID, distributorID)

	ape.Render(w, http.StatusCreated, responses.EmployeeInvites(invite))
}
