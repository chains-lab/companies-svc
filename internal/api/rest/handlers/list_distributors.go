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
)

func (s Service) ListDistributors(w http.ResponseWriter, r *http.Request) {
	filters := app.FilterDistributorList{}
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

	distributors, pag, err := s.app.ListDistributors(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to select distributors")

		switch {
		case errors.Is(err, errx.InvalidDistributorStatus):
			ape.RenderErr(w, problems.InvalidParameter("status", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	ape.Render(w, http.StatusOK, responses.DistributorCollection(distributors, pag))
}
