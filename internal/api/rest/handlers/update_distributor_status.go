package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) UpdateDistributorStatus(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get user from context")

		ape.RenderErr(w, problems.Unauthorized("user not found in context"))
		return
	}

	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor id: %s", chi.URLParam(r, "distributor_id"))

		ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
	}

	status := chi.URLParam(r, "status")
	if err := enum.ParseDistributorStatus(status); err != nil {
		s.Log(r).WithError(err).Errorf("invalid distributor status: %s", status)

		ape.RenderErr(w, problems.InvalidParameter("status", err))
		return
	}

	var distributor models.Distributor
	switch status {
	case enum.DistributorStatusActive:
		distributor, err = s.app.SetDistributorStatusActive(r.Context(), initiator.ID, distributorID)
		if err != nil {
			s.Log(r).WithError(err).Errorf("failed to set distributor %s status to active", distributorID)

			switch {
			default:
				ape.RenderErr(w, problems.InternalError())
			}
			return
		}

	case enum.DistributorStatusInactive:
		distributor, err = s.app.SetDistributorStatusInactive(r.Context(), initiator.ID, distributorID)
		if err != nil {
			s.Log(r).WithError(err).Errorf("failed to set distributor %s status to inactive", distributorID)

			switch {
			default:
				ape.RenderErr(w, problems.InternalError())
			}
			return
		}

	default:
		ape.RenderErr(w, problems.InvalidParameter("status", fmt.Errorf("unsupported status: %s", status)))
		return

	}

	s.Log(r).Infof("distributor %s status set to active successfully", distributor.ID)

	ape.Render(w, http.StatusOK, responses.Distributor(distributor))
}
