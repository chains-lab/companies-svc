package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	InitiatorID uuid.UUID
	Role        string
}

func (a App) SentInvite(ctx context.Context, params SentInviteParams) (models.Invite, error) {
	return a.employee.SentInvite(ctx, employee.SentInviteParams{
		InitiatorID: params.InitiatorID,
		Role:        params.Role,
	})
}
