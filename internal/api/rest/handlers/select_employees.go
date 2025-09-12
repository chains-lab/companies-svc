package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (s Service) ListEmployees(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := app.FilterEmployeeList{}

	if ids := q["distributor_id"]; len(ids) > 0 {
		filters.Distributors = make([]uuid.UUID, 0, len(ids))
		for _, raw := range ids {
			v, err := uuid.Parse(strings.TrimSpace(raw))
			if err != nil {
				s.Log(r).WithError(err).Errorf("invalid distributor ID format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
				return
			}
			filters.Distributors = append(filters.Distributors, v)
		}
	}

	if roles := q["role"]; len(roles) > 0 {
		filters.Roles = make([]string, 0, len(roles))
		for _, raw := range roles {
			filters.Roles = append(filters.Roles, raw)
		}
	}

	pagReq, sort := pagi.GetPagination(r)

	employees, pag, err := s.app.ListEmployees(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to select employees")

		switch {
		case errors.Is(err, errx.ErrorInvalidEmployeeRole):
			ape.RenderErr(w, problems.InvalidParameter("role", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeCollection(employees, pag))
}
