package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (a Service) CreateCompanyBlock(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateCompanyBlock(r)
	if err != nil {
		a.log.WithError(err).Error("failed to decode block company request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	//if req.Data.Attributes.CompanyID != chi.URLParam(r, "company_id") {
	//	ape.RenderErr(w,
	//		problems.InvalidParameter("company_id", fmt.Errorf("path ID and body ID do not match")),
	//		problems.InvalidPointer("/data/attributes.company_id", fmt.Errorf("path ID and body ID do not match")),
	//	)
	//
	//	return
	//}
	//
	//companyID, err := uuid.Parse(req.Data.Attributes.CompanyID)
	//if err != nil {
	//	a.log.WithError(err).Errorf("invalid company id: %s", req.Data.Attributes.CompanyID)
	//	ape.RenderErr(w, problems.BadRequest(validation.Errors{
	//		"company_id": err,
	//	})...)
	//
	//	return
	//}

	block, err := a.domain.company.CreteBlock(r.Context(), initiator.ID, req.Data.Attributes.CompanyId, req.Data.Attributes.Reason)
	if err != nil {
		a.log.WithError(err).Errorf("failed to block company")
		switch {
		case errors.Is(err, errx.ErrorcompanyHaveAlreadyActiveBlock):
			ape.RenderErr(w, problems.Conflict("company already have active block"))
		case errors.Is(err, errx.ErrorcompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("company %s blocked successfully by user %s", req.Data.Attributes.CompanyId, initiator.ID)

	responses.CompanyBlock(block)
}
