package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
)

func (s Service) RefuseMyEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	err = s.domain.employee.DeleteMe(r.Context(), initiator.ID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get employee")
		switch {
		case errors.Is(err, errx.ErrorOwnerCannotRefuseSelf):
			ape.RenderErr(w, problems.Forbidden("owner cannot refuse himself"))
		case errors.Is(err, errx.ErrorInitiatorIsNotEmployee):
			ape.RenderErr(w, problems.NotFound("employee not found"))
		case errors.Is(err, errx.ErrorCompanyNotFound):
			ape.RenderErr(w, problems.NotFound("company not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
