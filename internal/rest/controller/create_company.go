package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/company"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/requests"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (s Service) CreateCompany(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateCompany(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create company request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.company.CreateByEmployee(r.Context(), initiator.ID, company.CreateParams{
		Name: req.Data.Attributes.Name,
		Icon: req.Data.Attributes.Icon,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to create company")
		switch {
		case errors.Is(err, errx.ErrorCurrentEmployeeCannotCreateCompany):
			ape.RenderErr(w, problems.Conflict(
				fmt.Sprintf("current employee %s cannot create company", initiator.ID),
			))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("company %s created successfully", res.ID)

	ape.Render(w, http.StatusCreated, responses.Company(res))
}
