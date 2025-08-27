package handlers

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/pagi"
)

func (s Service) SelectDistributors(w http.ResponseWriter, r *http.Request) {
	filters := app.SelectDistributorsParams{}
	q := r.URL.Query()

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = make([]string, 0, len(sts))
		for _, raw := range sts {
			v := strings.TrimSpace(raw)
			if err := enum.ParseDistributorStatus(v); err != nil {
				s.Log(r).WithError(err).Errorf("invalid distributor status format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("status", err))
				return
			}
			filters.Statuses = append(filters.Statuses, v)
		}
	}

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &name
	}

	pagReq, sort := pagi.GetPagination(r)

	distributors, pag, err := s.app.SelectDistributors(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to select distributors")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	ape.Render(w, http.StatusOK, responses.DistributorCollection(distributors, pag))
}
