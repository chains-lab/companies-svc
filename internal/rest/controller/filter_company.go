package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/services/company"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/pagi"
)

func (s Service) FilterCompanies(w http.ResponseWriter, r *http.Request) {
	filters := company.FiltersParams{}
	q := r.URL.Query()

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = make([]string, 0, len(sts))
		for _, raw := range sts {
			filters.Statuses = append(filters.Statuses, strings.TrimSpace(raw))
		}
	}

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &name
	}

	pagReq, sort := pagi.GetPagination(r)

	companies, err := s.domain.company.Filter(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.log.WithError(err).Error("failed to select companies")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CompanyCollection(companies))
}
