package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (a Service) UpdateCompaniesStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCompanyStatus(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update company status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := a.domain.company.UpdateStatus(r.Context(), req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		a.log.WithError(err).Errorf("failed to set company %s status to active", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee has not enough rights"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsBlocked):
			ape.RenderErr(w, problems.Conflict("company is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("company %s status set to active successfully", res.ID)

	ape.Render(w, http.StatusOK, responses.Company(res))
}
