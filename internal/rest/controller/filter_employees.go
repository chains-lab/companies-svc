package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/pagi"
)

func (s Service) ListEmployees(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := employee.FilterParams{}

	if empls := q["employee_id"]; len(empls) > 0 {
		filters.EmployeeID = make([]uuid.UUID, 0, len(empls))
		for _, raw := range empls {
			id, err := uuid.Parse(raw)
			if err != nil {
				s.log.WithError(err).Errorf("invalid employee EmployeeID format")
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"employee_id": err,
				})...)
				return
			}

			filters.EmployeeID = append(filters.EmployeeID, id)
		}
	}

	if comps := q["company_id"]; len(comps) > 0 {
		filters.CompanyID = make([]uuid.UUID, 0, len(comps))
		for _, raw := range comps {
			id, err := uuid.Parse(raw)
			if err != nil {
				s.log.WithError(err).Errorf("invalid company EmployeeID format")
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"company_id": err,
				})...)
				return
			}

			filters.CompanyID = append(filters.CompanyID, id)
		}
	}

	if users := q["user_id"]; len(users) > 0 {
		filters.UserID = make([]uuid.UUID, 0, len(users))
		for _, raw := range users {
			id, err := uuid.Parse(raw)
			if err != nil {
				s.log.WithError(err).Errorf("invalid user EmployeeID format")
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"user_id": err,
				})...)
				return
			}

			filters.UserID = append(filters.UserID, id)
		}
	}

	if roles := q["role"]; len(roles) > 0 {
		filters.Roles = make([]string, 0, len(roles))
		for _, raw := range roles {
			filters.Roles = append(filters.Roles, raw)
		}
	}

	page, size := pagi.GetPagination(r)

	employees, err := s.domain.employee.Filter(r.Context(), filters, page, size)
	if err != nil {
		s.log.WithError(err).Error("failed to select employees")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeCollection(employees))
}
