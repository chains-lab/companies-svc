package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid user ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"user_id": err,
		})...)

		return
	}

	companyID, err := uuid.Parse(chi.URLParam(r, "company_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid company ID format")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"company_id": err,
		})...)

		return
	}

	err = s.domain.employee.DeleteByInitiatorID(r.Context(), initiator.ID, userID, companyID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to delete employee with user_id: %s", userID)
		switch {
		case errors.Is(err, errx.ErrorEmployeeNotFound):
			ape.RenderErr(w, problems.NotFound("employee not found"))
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.Forbidden("initiator is not an employee"))
		case errors.Is(err, errx.ErrorCannotDeleteYourself):
			ape.RenderErr(w, problems.Forbidden("cannot delete yourself"))
		case errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("initiator employee have not enough rights"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		case errors.Is(err, errx.ErrorCompanyIsNotActive):
			ape.RenderErr(w, problems.Forbidden("company is not active"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	s.log.Infof("employee %s deleted successfully", userID)

	w.WriteHeader(http.StatusNoContent)
	return
}
