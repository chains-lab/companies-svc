package distributor

import (
	"context"
	"database/sql"

	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/google/uuid"
)

type distributorsQ interface {
	New() dbx.DistributorsQ
	Insert(ctx context.Context, input dbx.Distributor) error
	Get(ctx context.Context) (dbx.Distributor, error)
	Select(ctx context.Context) ([]dbx.Distributor, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.DistributorsQ
	FilterStatus(status ...string) dbx.DistributorsQ
	LikeName(name string) dbx.DistributorsQ

	OrderByName(ascend bool) dbx.DistributorsQ

	Page(limit, offset uint64) dbx.DistributorsQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

type blockagesQ interface {
	New() dbx.BlockQ

	Insert(ctx context.Context, input dbx.Blockages) error
	Get(ctx context.Context) (dbx.Blockages, error)
	Select(ctx context.Context) ([]dbx.Blockages, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.BlockQ
	FilterDistributorID(distributorID ...uuid.UUID) dbx.BlockQ
	FilterInitiatorID(initiatorID ...uuid.UUID) dbx.BlockQ
	FilterStatus(status ...string) dbx.BlockQ

	OrderByBlockedAt(ascending bool) dbx.BlockQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.BlockQ
}

type Distributor struct {
	distributor distributorsQ
	block       blockagesQ
}

func NewDistributor(db *sql.DB) Distributor {
	return Distributor{
		distributor: dbx.NewDistributorsQ(db),
		block:       dbx.NewBlockagesQ(db),
	}
}
