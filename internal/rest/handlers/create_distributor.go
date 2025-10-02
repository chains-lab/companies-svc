package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/rest/requests"
	"github.com/chains-lab/distributors-svc/internal/rest/responses"
)

func (a Adapter) CreateDistributor(w http.ResponseWriter, r *http.Request) {
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

	distributor, err := a.app.CreateDistributor(
		r.Context(),
		initiator.ID,
		req.Data.Attributes.Name,
		req.Data.Attributes.Icon,
	)
	if err != nil {
		a.log.WithError(err).Errorf("failed to create distributor")
		switch {
		case errors.Is(err, errx.ErrorCurrentEmployeeCannotCreateDistributor):
			ape.RenderErr(w, problems.PreconditionFailed("Current employee can not create distributor"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("distributor %s created successfully", distributor.ID)

	ape.Render(w, http.StatusCreated, responses.Distributor(distributor))
}
