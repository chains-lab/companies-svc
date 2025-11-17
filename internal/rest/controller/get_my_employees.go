package controller

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (s Service) GetMyEmployees(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	page, size := pagi.GetPagination(r)

	emps, err := s.domain.employee.Filter(r.Context(),
		employee.FilterParams{
			UserID: []uuid.UUID{initiator.ID},
		},
		page, size,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get employee")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeCollection(emps))
}
