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
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.UpdateEmployee(r)
	if err != nil {
		s.log.WithError(err).Error("failed to parse update employee request")

		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	params := employee.UpdateParams{
		Role:     req.Data.Attributes.Role,
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	}

	res, err := s.domain.employee.UpdateByEmployee(r.Context(), req.Data.Id, initiator.ID, params)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update employee for EmployeeID: %s", req.Data.Id)
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRight):
			ape.RenderErr(w, problems.Forbidden("initiator have not enough rights"))
		case errors.Is(err, errx.ErrorEmployeeNotFound):
			ape.RenderErr(w, problems.NotFound("employee not found"))
		case errors.Is(err, errx.ErrorInvalidEmployeeRole):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"data/attributes/role": err,
			})...)
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsNotActive):
			ape.RenderErr(w, problems.Forbidden("company is not active"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		s.log.WithError(err).Errorf("internal error when updating employee for EmployeeID: %s", req.Data.Id)
		return
	}

	ape.Render(w, http.StatusOK, responses.Employee(res))
}
