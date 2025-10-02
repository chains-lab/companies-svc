package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetEmployee(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid user ID format")
		ape.RenderErr(w, problems.InvalidParameter("user_id", err))

		return
	}

	employee, err := a.app.GetEmployee(r.Context(), userID)
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

	ape.Render(w, http.StatusOK, responses.Employee(employee))
}
