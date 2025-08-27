package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) InteractToInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
		return
	}

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid invite ID format: %s", chi.URLParam(r, "invite_id"))

		ape.RenderErr(w, problems.InvalidParameter("invite_id", err))
		return
	}

	status := chi.URLParam(r, "status")
	err = enum.ParseInviteStatus(chi.URLParam(r, "status"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid invite status format: %s", chi.URLParam(r, "status"))

		ape.RenderErr(w, problems.InvalidParameter("status", err))
		return
	}

	var invite models.EmployeeInvite
	switch status {
	case enum.InviteStatusAccepted:
		invite, err = s.app.AcceptInvite(r.Context(), initiator.ID, inviteID)
		if err != nil {
			s.Log(r).WithError(err).Errorf("failed to accept invite")

			switch {
			case errors.Is(err, errx.InviteNotFound):
				ape.RenderErr(w, problems.NotFound("Invite not found"))
			case errors.Is(err, errx.InviteIsNotActive):
				ape.RenderErr(w, problems.Conflict("Invite is not active"))
			case errors.Is(err, errx.InviteIsNotForInitiator):
				ape.RenderErr(w, problems.Forbidden("Initiator is not for initiator"))
			case errors.Is(err, errx.InitiatorIsAlreadyEmployee):
				ape.RenderErr(w, problems.Conflict("Initiator is already employee"))
			default:
				ape.RenderErr(w, problems.InternalError())
			}
			return
		}
	case enum.InviteStatusRejected:
		invite, err = s.app.RejectInvite(r.Context(), initiator.ID, inviteID)
		if err != nil {
			s.Log(r).WithError(err).Errorf("failed to reject invite")

			switch {
			case errors.Is(err, errx.InviteNotFound):
				ape.RenderErr(w, problems.NotFound("Invite not found"))
			case errors.Is(err, errx.InviteIsNotActive):
				ape.RenderErr(w, problems.Conflict("Invite is not active"))
			case errors.Is(err, errx.InviteIsNotForInitiator):
				ape.RenderErr(w, problems.Forbidden("Initiator is not for initiator"))
			default:
				ape.RenderErr(w, problems.InternalError())
			}
			return
		}
	case enum.InviteStatusWithdrawn:
		invite, err = s.app.WithdrawInvite(r.Context(), initiator.ID, inviteID)
		if err != nil {
			s.Log(r).WithError(err).Errorf("failed to withdraw invite")

			switch {
			case errors.Is(err, errx.InviteNotFound):
				ape.RenderErr(w, problems.NotFound("Invite not found"))
			case errors.Is(err, errx.InviteIsNotActive):
				ape.RenderErr(w, problems.Conflict("Invite is not active"))
			case errors.Is(err, errx.InitiatorEmployeeHaveNotEnoughRights):
				ape.RenderErr(w, problems.Forbidden("Initiator employee have not enough rights"))
			case errors.Is(err, errx.InviteIsNotForInitiator):
				ape.RenderErr(w, problems.Forbidden("Initiator is not for initiator"))
			default:
				ape.RenderErr(w, problems.InternalError())
			}
			return
		}
	default:
		s.Log(r).WithError(err).Errorf("invalid invite status for update: %s", status)

		ape.RenderErr(w, problems.InvalidParameter("status",
			fmt.Errorf("invalid status: %s muste be one of: [%s, %s, %s]",
				status, enum.InviteStatusAccepted, enum.InviteStatusRejected, enum.InviteStatusWithdrawn),
		),
		)
		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeInvites(invite))
	return
}
