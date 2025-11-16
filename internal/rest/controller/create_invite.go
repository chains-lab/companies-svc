package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/invite"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (s Service) CreateInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.CreateInvite(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create employee invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.invite.Create(r.Context(), initiator.ID, invite.CreateParams{
		CompanyID: req.Data.Attributes.CompanyId,
		UserID:    req.Data.Attributes.UserId,
		Role:      req.Data.Attributes.Role,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to create employee invite")
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not employee"))
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployeeInThisCompany):
			ape.RenderErr(w, problems.Forbidden("initiator is not employee of this company"))
		case errors.Is(err, errx.ErrorUserAlreadyEmployee):
			ape.RenderErr(w, problems.Conflict("user is already employee"))
		case errors.Is(err, errx.ErrorNotEnoughRight):
			ape.RenderErr(w, problems.Forbidden("initiator have not enough rights"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsNotActive):
			ape.RenderErr(w, problems.Conflict("company is not active"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Invites(res))
}
