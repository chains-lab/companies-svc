package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/distributors-svc/internal/api/rest/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	blockID, err := uuid.Parse(chi.URLParam(r, "block_id"))
	if err != nil {
		s.Log(r).WithError(err).Errorf("invalid block ID format")

		ape.RenderErr(w, problems.InvalidParameter("block_id", err))
		return
	}

	block, err := s.app.GetBlock(r.Context(), blockID)
	if err != nil {
		s.Log(r).WithError(err).Errorf("failed to get block, ID: %s", blockID)

		switch {
		case errors.Is(err, errx.DistributorBlockNotFound):
			ape.RenderErr(w, problems.NotFound("Block not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	ape.Render(w, http.StatusOK, responses.DistributorBlock(block))
	return
}
