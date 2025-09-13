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

func (a Adapter) ListBlockages(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := app.FilterBlockagesList{}

	if ids := q["distributor_id"]; len(ids) > 0 {
		filters.Distributors = make([]uuid.UUID, 0, len(ids))
		for _, raw := range ids {
			v, err := uuid.Parse(strings.TrimSpace(raw))
			if err != nil {
				a.log.WithError(err).Errorf("invalid distributor ID format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
				return
			}
			filters.Distributors = append(filters.Distributors, v)
		}
	}

	if ids := q["initiator_id"]; len(ids) > 0 {
		filters.Initiators = make([]uuid.UUID, 0, len(ids))
		for _, raw := range ids {
			v, err := uuid.Parse(strings.TrimSpace(raw))
			if err != nil {
				a.log.WithError(err).Errorf("invalid initator ID format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("initiator_id", err))
				return
			}
			filters.Initiators = append(filters.Initiators, v)
		}
	}

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = make([]string, 0, len(sts))
		for _, raw := range sts {
			filters.Statuses = append(filters.Statuses, strings.TrimSpace(raw))
		}
	}

	pagReq, sort := pagi.GetPagination(r)

	blocks, pag, err := a.app.ListBlockages(r.Context(), filters, pagReq, sort)
	if err != nil {
		a.log.WithError(err).Error("failed to select distributor blocks")
		switch {
		case errors.Is(err, errx.ErrorInvalidDistributorBlockStatus):
			ape.RenderErr(w, problems.InvalidParameter("status", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlockCollection(blocks, pag))
}
