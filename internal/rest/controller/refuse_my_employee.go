package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
)

func (a Service) RefuseMyEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	err = a.domain.employee.RefuseMe(r.Context(), initiator.ID)
	if err != nil {
		a.log.WithError(err).Errorf("failed to get employee")
		switch {
		case errors.Is(err, errx.ErrorEmployeeNotFound):
			ape.RenderErr(w, problems.Unauthorized("employee not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
