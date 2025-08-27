package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
)

func (s Service) UpdateEmployeeRole(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		s.Log(r).WithError(err).Error("failed to get userID from context")

		ape.RenderErr(w, problems.Unauthorized("user not found in context"))
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid userID id: %s", chi.URLParam(r, "user_id"))

		ape.RenderErr(w, problems.InvalidParameter("user_id", err))
		return
	}

	role := chi.URLParam(r, "role")
	if role == "" {
		s.Log(r).Error("role parameter is required")

		ape.RenderErr(w, problems.InvalidParameter("role", fmt.Errorf("role parameter is required")))
		return
	}

	err = enum.ParseDistributorStatus(role)
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid role: %s", role)

		ape.RenderErr(w, problems.InvalidParameter("role", err))
		return
	}

	res, err := s.app.UpdateEmployeeRole(r.Context(), initiator.ID, userID, role)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to update employee role for ID: %s", userID)

		switch {

		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Employee(res))
}
