package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Service) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	token := chi.URLParam(r, "token")

	invite, err := a.domain.employee.AcceptInvite(r.Context(), initiator.ID, token)
	if err != nil {
		a.log.WithError(err).Error("failed to answer to invite")
		switch {
		case errors.Is(err, errx.ErrorInvalidInviteToken):
			ape.RenderErr(w, problems.Unauthorized("invalid invite token"))
		case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
			ape.RenderErr(w, problems.Conflict("invite already answered"))
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteExpired):
			ape.RenderErr(w, problems.Conflict("invite expired"))
		case errors.Is(err, errx.ErrorUnexpectedInviteStatus):
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invites(invite))
}
