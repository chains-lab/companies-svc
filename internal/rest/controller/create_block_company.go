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
)

func (s Service) CreateCompanyBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateCompanyBlock(r)
	if err != nil {
		s.log.WithError(err).Error("failed to decode block company request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	block, err := s.domain.block.Crete(r.Context(), initiator.ID, req.Data.Attributes.CompanyId, req.Data.Attributes.Reason)
	if err != nil {
		s.log.WithError(err).Errorf("failed to block company")
		switch {
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyHaveAlreadyActiveBlock):
			ape.RenderErr(w, problems.Conflict("company already have active block"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("company %s blocked successfully by user %s", req.Data.Attributes.CompanyId, initiator.ID)

	responses.CompanyBlock(block)
}
