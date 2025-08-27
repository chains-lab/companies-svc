package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/distributors-svc/internal/api/rest/requests"
)

func (s Service) UpdateDistributor(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		return
	}

	req, err := requests.UpdateDistributor(r)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to parse update distributor request")

		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	if req.Data.Id != chi.URLParam(r, "distributor_id") {
		ape.RenderErr(w, problems.InvalidParameter("distributor_id", fmt.Errorf("path ID and body ID do not match")))

		return
	}

	distributorID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor id: %s", req.Data.Id)

		ape.RenderErr(w, problems.InvalidParameter("id", err))
		return
	}

	input := app.UpdateDistributorInput{}
	if req.Data.Attributes.Name != nil {
		input.Name = req.Data.Attributes.Name
	}
	if req.Data.Attributes.Icon != nil {
		input.Icon = req.Data.Attributes.Icon
	}

	distributor, err := s.app.UpdateDistributor(r.Context(), initiator.ID, distributorID, input)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to update distributor name for ID: %s", distributorID)

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("distributor %s updated successfully", distributorID)

	ape.Render(w, http.StatusOK, responses.Distributor(distributor))
}
