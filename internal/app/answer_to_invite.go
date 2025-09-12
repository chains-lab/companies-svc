package app

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/google/uuid"
)

type AnswerToInviteParams struct {
	Token  string
	Answer string
}

func (a App) AnswerToInvite(ctx context.Context, userID uuid.UUID, params AnswerToInviteParams) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.employee.AnsweredInvite(ctx, userID, employee.AnsweredInviteParams{
			Token:  params.Token,
			Status: params.Answer,
		})
		return err
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, nil
}
