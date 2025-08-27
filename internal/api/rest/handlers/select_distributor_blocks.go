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
	"github.com/google/uuid"
)

func (s Service) SelectDistributorBlocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := app.SelectBlockagesParams{}

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

	if ids := q["initiator_id"]; len(ids) > 0 {
		filters.Initiators = make([]uuid.UUID, 0, len(ids))
		for _, raw := range ids {
			v, err := uuid.Parse(strings.TrimSpace(raw))
			if err != nil {
				s.Log(r).WithError(err).Errorf("invalid initator ID format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("initiator_id", err))
				return
			}
			filters.Initiators = append(filters.Initiators, v)
		}
	}

	if sts := q["status"]; len(sts) > 0 {
		filters.Statuses = make([]string, 0, len(sts))
		for _, raw := range sts {
			v := strings.TrimSpace(raw)
			if err := enum.ParseBlockStatus(v); err != nil {
				s.Log(r).WithError(err).Errorf("invalid block status format: %s", raw)
				ape.RenderErr(w, problems.InvalidParameter("status", err))
				return
			}
			filters.Statuses = append(filters.Statuses, v)
		}
	}

	pagReq, sort := pagi.GetPagination(r)

	blocks, pag, err := s.app.SelectBlockages(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to select distributor blocks")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlockCollection(blocks, pag))
}
