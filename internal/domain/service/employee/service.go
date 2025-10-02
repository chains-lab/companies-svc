package employee

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/infra/jwtmanager"
	"github.com/google/uuid"
)

type Service struct {
	db  database
	jwt jwtmanager.Manager
}

func NewService(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
