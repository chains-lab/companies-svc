package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/company"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (s Service) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.UpdateCompany(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update company request")

		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	input := company.UpdateParams{
		Name: req.Data.Attributes.Name,
		Icon: req.Data.Attributes.Icon,
	}

	res, err := s.domain.company.UpdateByInitiator(r.Context(), initiator.ID, req.Data.Id, input)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update company name for ID: %s", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee has not enough rights"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsNotActive):
			ape.RenderErr(w, problems.Forbidden("company is not active"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("company %s updated successfully", fmt.Sprint(res.ID))

	ape.Render(w, http.StatusOK, responses.Company(res))
}
