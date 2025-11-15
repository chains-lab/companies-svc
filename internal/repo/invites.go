package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (r *Repo) CreateInvite(ctx context.Context, input models.Invite) error {
	return r.sql.invites.New().Insert(ctx, inviteModelToSchema(input))
}

func (r *Repo) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	row, err := r.sql.invites.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Invite{}, nil
	case err != nil:
		return models.Invite{}, err
	}

	return inviteSchemaToModel(row), nil
}

func (r *Repo) UpdateInviteStatus(ctx context.Context, ID uuid.UUID, status string) error {
	return r.sql.invites.New().
		FilterID(ID).
		UpdateStatus(status).
		Update(ctx)
}
