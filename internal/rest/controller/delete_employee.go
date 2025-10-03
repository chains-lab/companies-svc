package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (a Service) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid user ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	distributorID, err := uuid.Parse(chi.URLParam(r, "distributor_id"))
	if err != nil {
		a.log.WithError(err).Errorf("invalid distributor ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"distributor_id": err,
		})...)

		return
	}

	err = a.domain.employee.Delete(r.Context(), initiator.ID, userID, distributorID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to delete employee with user_id: %s", userID)
		switch {
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee have not enough rights"))
		case errors.Is(err, errx.ErrorDistributorNotFound):
			ape.RenderErr(w, problems.NotFound("distributor not found"))
		case errors.Is(err, errx.ErrorEmployeeNotFound):
			ape.RenderErr(w, problems.NotFound("employee not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("employee %s deleted successfully", userID)

	w.WriteHeader(http.StatusNoContent)
	return
}
