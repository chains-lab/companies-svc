package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
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

	distributor, err := s.app.SetDistributorStatus(r.Context(), initiator.ID, distributorID, status)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to set distributor %s status to active", distributorID)

		switch {
		case errors.Is(err, errx.ErrorInitiatorNotEmployee):
			ape.RenderErr(w, problems.PreconditionFailed("initiator is not an employee"))
		case errors.Is(err, errx.ErrorInitiatorEmployeeHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator have not enough rights"))
		case errors.Is(err, errx.UnexpectedDistributorSetStatus):
			ape.RenderErr(w, problems.PreconditionFailed("distributor status is blocked"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	s.Log(r).Infof("distributor %s status set to active successfully", distributor.ID)

	ape.Render(w, http.StatusOK, responses.Distributor(distributor))
}
