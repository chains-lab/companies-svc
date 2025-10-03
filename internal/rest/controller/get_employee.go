package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) GetEmployee(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid user ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	emp, err := a.domain.employee.Get(r.Context(), employee.GetParams{
		UserID: &userID,
	})
	if err != nil {
		a.log.WithError(err).Errorf("failed to get employee")
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
