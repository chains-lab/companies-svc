package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/rest/meta"
	"github.com/chains-lab/companies-svc/internal/rest/responses"
)

func (a Service) GetMyEmployee(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	me, err := a.domain.employee.Get(r.Context(), employee.GetParams{
		UserID: &initiator.ID,
	})
	if err != nil {
		a.log.WithError(err).Errorf("failed to get employee")
		switch {
		case errors.Is(err, errx.ErrorOwnerCannotRefuseSelf):
			ape.RenderErr(w, problems.Forbidden("owner cannot refuse self"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Employee(me))
}
