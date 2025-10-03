package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d *Database) CreateInvite(ctx context.Context, input models.Invite) error {
	schema := inviteModelToSchema(input)

	return d.sql.invites.New().Insert(ctx, schema)
}

func (d *Database) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	row, err := d.sql.invites.New().FilterID(ID).Get(ctx)
	if err != nil {
		return models.Invite{}, err
	}

	return inviteSchemaToModel(row), nil
}

func (d *Database) UpdateInviteStatus(ctx context.Context, ID, UserID uuid.UUID, status string, acceptedAt time.Time) error {
	return d.sql.invites.New().
		FilterID(ID).
		UpdateStatus(status).
		UpdateAnsweredAt(acceptedAt).
		UpdateUserID(UserID).
		Update(ctx)
}

func inviteModelToSchema(m models.Invite) pgdb.Invite {
	res := pgdb.Invite{
		ID:            m.ID,
		Status:        m.Status,
		Role:          m.Role,
		DistributorID: m.DistributorID,
		Token:         m.Token,
		CreatedAt:     m.CreatedAt,
		ExpiresAt:     m.ExpiresAt,
	}
	if m.UserID != nil {
		res.UserID = uuid.NullUUID{
			UUID:  *m.UserID,
			Valid: true,
		}
	}
	if m.AnsweredAt != nil {
		res.AnsweredAt = sql.NullTime{
			Time:  *m.AnsweredAt,
			Valid: true,
		}
	}

	return res
}

func inviteSchemaToModel(m pgdb.Invite) models.Invite {
	res := models.Invite{
		ID:            m.ID,
		Status:        m.Status,
		Role:          m.Role,
		DistributorID: m.DistributorID,
		Token:         m.Token,
		CreatedAt:     m.CreatedAt,
		ExpiresAt:     m.ExpiresAt,
	}
	if m.UserID.Valid {
		res.UserID = &m.UserID.UUID
	}
	if m.AnsweredAt.Valid {
		res.AnsweredAt = &m.AnsweredAt.Time
	}

	return res
}
