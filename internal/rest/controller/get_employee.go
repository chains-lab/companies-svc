package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) GetEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	cityID, err := uuid.Parse(chi.URLParam(r, "employee_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid employee EmployeeID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"employee_id": err,
		})...)

		return
	}

	emp, err := s.domain.employee.Get(r.Context(), initiator.ID, cityID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get employee")
		switch {
		case errors.Is(err, errx.ErrorEmployeeNotFound):
			ape.RenderErr(w, problems.NotFound("Employee not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Employee(emp))
}
