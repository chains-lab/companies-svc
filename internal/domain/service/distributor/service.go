package distributor

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{
		db: db,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmployee(ctx context.Context, input models.Employee) error
	GetUserEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error)

	CreateDistributor(ctx context.Context, input models.Distributor) (models.Distributor, error)
	GetDistributorByID(ctx context.Context, ID uuid.UUID) (models.Distributor, error)
	FilterDistributors(ctx context.Context, filters Filters, page, size uint64) (models.DistributorCollection, error)
	UpdateDistributor(ctx context.Context, ID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	UpdateDistributorStatus(ctx context.Context, ID uuid.UUID, status string, updatedAt time.Time) error

	CreateDistributorBlock(ctx context.Context, input models.DistributorBlock) error
	GetDistributorBlockByID(ctx context.Context, ID uuid.UUID) (models.DistributorBlock, error)
	GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.DistributorBlock, error)
	FilterDistributorBlocks(ctx context.Context, filters FilterBlockages, page, size uint64) (models.DistributorBlockCollection, error)
	CancelActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID, canceledAt time.Time) error
}
