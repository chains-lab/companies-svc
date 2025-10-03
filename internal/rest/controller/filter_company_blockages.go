package controller

import (
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/service/block"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
	"github.com/chains-lab/pagi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) FilterBlockages(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := block.FilterParams{}

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
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		err := enum.CheckCompanyStatus(v)
		if err != nil {
			a.log.WithError(err).Errorf("invalid company block status format")
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

	blocks, err := a.domain.block.Filter(r.Context(), filters, page, size)
	if err != nil {
		a.log.WithError(err).Error("failed to select company blocks")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.CompanyBlockCollection(blocks))
}
