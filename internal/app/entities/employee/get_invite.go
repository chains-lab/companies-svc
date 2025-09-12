package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/jwtmanager"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
)

func (e Employee) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	inv, err := e.invite.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.Invite{}, errx.ErrorInternal.Raise(fmt.Errorf("get invite by ID, cause %w", err))
		}
	}

	token, err := e.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		ID:            inv.ID,
		DistributorID: inv.DistributorID,
		Role:          inv.Role,
		ExpiredAt:     inv.ExpiresAt,
		CreatedAt:     inv.CreatedAt,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	return inviteFromDB(inv, token), nil
}
