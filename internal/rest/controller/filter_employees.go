package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"

	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/chains-lab/pagi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) ListEmployees(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := employee.FilterParams{}

	if v := strings.TrimSpace(q.Get("company_id")); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			a.log.WithError(err).Errorf("invalid company ID format")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"company_id": err,
			})...)

			return
		}
		filters.CompanyID = &id
	}

	if roles := q["role"]; len(roles) > 0 {
		filters.Roles = make([]string, 0, len(roles))
		for _, raw := range roles {
			filters.Roles = append(filters.Roles, raw)
		}
	}

	page, size := pagi.GetPagination(r)

	employees, err := a.domain.employee.Filter(r.Context(), filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to select employees")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeCollection(employees))
}
