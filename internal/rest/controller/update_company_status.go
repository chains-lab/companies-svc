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

func (s Service) UpdateCompaniesStatus(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.UpdateCompanyStatus(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update company status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.company.UpdateStatusByInitiator(r.Context(), initiator.ID, req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		s.log.WithError(err).Errorf("failed to set company %s status to active", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee has not enough rights"))
		case errors.Is(err, errx.ErrorInvalidCompanyBlockStatus):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/status": err,
			})...)
		case errors.Is(err, errx.ErrorCannotSetCompanyStatusBlocked):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/status": err,
			})...)
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsBlocked):
			ape.RenderErr(w, problems.Conflict("company is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("company %s status set to active successfully", res.ID)

	ape.Render(w, http.StatusOK, responses.Company(res))
}
