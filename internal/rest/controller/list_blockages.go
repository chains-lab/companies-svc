package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) ListBlockages(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := distributor.FilterBlockages{}

	if v := strings.TrimSpace(q.Get("distributor_id")); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			a.log.WithError(err).Errorf("invalid distributor ID format")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"distributor_id": err,
			})...)
			return
		}
		filters.DistributorID = &id
	}
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		err := enum.CheckDistributorBlockStatus(v)
		if err != nil {
			a.log.WithError(err).Errorf("invalid distributor block status format")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"status": err,
			})...)
			return
		}
		filters.Status = &v
	}
	if v := strings.TrimSpace(q.Get("initiator_id")); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			a.log.WithError(err).Errorf("invalid initiator ID format")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"initiator_id": err,
			})...)
			return
		}
		filters.InitiatorID = &id
	}

	// Pagination

	page, size := pagi.GetPagination(r)

	blocks, err := a.domain.distributor.FilterBlockages(r.Context(), filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to select distributor blocks")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlockCollection(blocks))
}
