package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid invite ID format")

		ape.RenderErr(w, problems.InvalidParameter("invite_id", err))
		return
	}

	invite, err := s.app.GetInvite(r.Context(), inviteID)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get invite with ID: %s", inviteID)

		switch {
		case errors.Is(err, errx.InviteNotFound):
			ape.RenderErr(w, problems.NotFound("Invite not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeInvites(invite))
	return
}
