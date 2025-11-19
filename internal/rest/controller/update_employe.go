package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
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

	companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid company EmployeeID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"company_id": err,
		})...)

		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid user EmployeeID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	ids := strings.Split(req.Data.Id, ":")
	if len(ids) != 2 {
		s.log.Error("invalid employee id format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("invalid id: %s, need format uuid:uuid look like user_id:company_id", req.Data.Id),
		})...)

		return
	}
	if ids[0] != userID.String() {
		s.log.Error("employee id user_id does not match url user_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("employee id user_id does not match url user_id"),
		})...)

		return
	}
	if ids[1] != companyID.String() {
		s.log.Error("employee id company_id does not match url company_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf("employee id company_id does not match url company_id"),
		})...)

		return
	}

	params := employee.UpdateParams{
		Role:     req.Data.Attributes.Role,
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	}

	res, err := s.domain.employee.UpdateByEmployee(r.Context(), initiator.ID, userID, companyID, params)
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
