package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d *Database) CreateInvite(ctx context.Context, input models.Invite) error {
	return d.sql.invites.New().Insert(ctx, inviteModelToSchema(input))
}

func (d *Database) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	row, err := d.sql.invites.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Invite{}, nil
	case err != nil:
		return models.Invite{}, err
	}

	return inviteSchemaToModel(row), nil
}

func (d *Database) UpdateInviteStatus(
	ctx context.Context,
	ID, UserID uuid.UUID,
	status string,
	acceptedAt time.Time,
) error {
	return d.sql.invites.New().
		FilterID(ID).
		UpdateStatus(status).
		UpdateAnsweredAt(acceptedAt).
		UpdateUserID(UserID).
		Update(ctx)
}
