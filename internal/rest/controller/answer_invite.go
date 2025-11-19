package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) ReplyInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.ReplyInvite(r)
	if err != nil {
		s.log.WithError(err).Error("invalid reply invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	invite, err := s.domain.invite.Reply(r.Context(), initiator.ID, req.Data.Id, req.Data.Attributes.Reply)
	if err != nil {
		s.log.WithError(err).Error("failed to reply to invite")
		switch {
		case errors.Is(err, errx.ErrorInvalidInviteStatus):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/reply": err,
			})...)

		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteExpired):
			ape.RenderErr(w, problems.Conflict("invite expired"))
		case errors.Is(err, errx.ErrorInviteAlreadyReplyed):
			ape.RenderErr(w, problems.Conflict("invite already replyed"))
		case errors.Is(err, errx.ErrorInviteNotForThisUser):
			ape.RenderErr(w, problems.Unauthorized("invite not for user"))
		case errors.Is(err, errx.ErrorCompanyIsNotActive):
			ape.RenderErr(w, problems.Forbidden("company is not active"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invites(invite))
}
