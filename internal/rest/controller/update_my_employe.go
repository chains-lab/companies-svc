package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (s Service) UpdateMyEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.UpdateMyEmployee(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update employee request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	params := employee.UpdateMyParams{
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	}

	res, err := s.domain.employee.UpdateMy(r.Context(), initiator.ID, params)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update employee for ID: %s", initiator.ID)
		switch {
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		s.log.WithError(err).Errorf("internal error when updating employee for ID: %s", initiator.ID)
		return
	}

	ape.Render(w, http.StatusOK, responses.Employee(res))
}
