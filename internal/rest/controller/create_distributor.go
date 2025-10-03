package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
)

func (a Service) CreateDistributor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateDistributor(r)
	if err != nil {
		a.log.WithError(err).Errorf("invalid create distributor request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := a.domain.distributor.Create(r.Context(), distributor.CreateParams{
		InitiatorID: initiator.ID,
		Name:        req.Data.Attributes.Name,
		Icon:        req.Data.Attributes.Icon,
	})
	if err != nil {
		a.log.WithError(err).Errorf("failed to create distributor")
		switch {
		case errors.Is(err, errx.ErrorCurrentEmployeeCannotCreateDistributor):
			ape.RenderErr(w, problems.Conflict(
				fmt.Sprintf("current employee %s cannot create distributor", initiator.ID),
			))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor %s created successfully", res.ID)

	ape.Render(w, http.StatusCreated, responses.Distributor(res))
}
