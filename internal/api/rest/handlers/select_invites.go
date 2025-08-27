package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (s Service) SelectInvites(w http.ResponseWriter, r *http.Request) {
	filters := app.SelectInvitesParams{}
	q := r.URL.Query()

	if ids := q["distributor_id"]; len(ids) > 0 {
		for _, idStr := range ids {
			id, err := uuid.Parse(idStr)
			if err != nil {
				s.Log(r).WithError(err).Errorf("invalid distributor ID format: %s", idStr)
				ape.RenderErr(w, problems.InvalidParameter("distributor_id", err))
				return
			}
			filters.Distributors = append(filters.Distributors, id)
		}
	}

	if ids := q["user_id"]; len(ids) > 0 {
		for _, idStr := range ids {
			id, err := uuid.Parse(idStr)
			if err != nil {
				s.Log(r).WithError(err).Errorf("invalid user ID format: %s", idStr)
				ape.RenderErr(w, problems.InvalidParameter("user_id", err))
				return
			}
			filters.ForUsers = append(filters.ForUsers, id)
		}
	}

	if ids := q["invited_by"]; len(ids) > 0 {
		for _, idStr := range ids {
			id, err := uuid.Parse(idStr)
			if err != nil {
				s.Log(r).WithError(err).Errorf("invalid invited_by format: %s", idStr)
				ape.RenderErr(w, problems.InvalidParameter("invited_by", err))
				return
			}
			filters.Inviters = append(filters.Inviters, id)
		}
	}

	if sts := q["status"]; len(sts) > 0 {
		for _, st := range sts {
			if err := enum.ParseInviteStatus(st); err != nil {
				s.Log(r).WithError(err).Errorf("invalid invite status: %s", st)
				ape.RenderErr(w, problems.InvalidParameter("status", err))
				return
			}
			filters.Statuses = append(filters.Statuses, st)
		}
	}

	if roles := q["role"]; len(roles) > 0 {
		for _, role := range roles {
			if err := enum.ParseEmployeeRole(role); err != nil {
				s.Log(r).WithError(err).Errorf("invalid role: %s", role)
				ape.RenderErr(w, problems.InvalidParameter("role", err))
				return
			}
			filters.Roles = append(filters.Roles, role)
		}
	}

	pagReq, sort := pagi.GetPagination(r)

	invites, pag, err := s.app.SelectInvites(r.Context(), filters, pagReq, sort)
	if err != nil {
		s.Log(r).WithError(err).Error("failed to select invites")

		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.EmployeeInvitesCollection(invites, pag))
}
