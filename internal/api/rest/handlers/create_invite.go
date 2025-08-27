package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/google/uuid"
)

func (s Service) CreateInvite(w http.ResponseWriter, r *http.Request) {
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

	role := req.Data.Attributes.Role
	if err := enum.ParseDistributorStatus(role); err != nil {
		s.Log(r).WithError(err).Errorf("invalid role: %s", role)

		ape.RenderErr(w, problems.InvalidParameter("role", err))
		return
	}

	s.Log(r).Infof("invite sent successfully for user ID: %s in distributor ID: %s", userID, distributorID)

	invite, err := s.app.CreateEmployeeInvite(r.Context(), initiator.ID, userID, distributorID, role)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to create employee invite")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.EmployeeInvites(invite))
}
