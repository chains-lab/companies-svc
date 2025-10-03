package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (a Service) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateCompany(r)
	if err != nil {
		a.log.WithError(err).Error("failed to parse update company request")

		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	input := company.UpdateParams{}
	if req.Data.Attributes.Name != nil {
		input.Name = req.Data.Attributes.Name
	}
	if req.Data.Attributes.Icon != nil {
		input.Icon = req.Data.Attributes.Icon
	}

	res, err := a.domain.company.Update(r.Context(), req.Data.Id, input)
	if err != nil {
		a.log.WithError(err).Errorf("failed to update company name for ID: %s", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee has not enough rights"))
		case errors.Is(err, errx.ErrorcompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorcompanyIsBlocked):
			ape.RenderErr(w, problems.Conflict("company is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("company %s updated successfully", fmt.Sprint(res.ID))

	ape.Render(w, http.StatusOK, responses.Company(res))
}
