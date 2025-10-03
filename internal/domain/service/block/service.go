package block

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
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

	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)
	UpdateCompaniesStatus(ctx context.Context, ID uuid.UUID, status string, updatedAt time.Time) error

	CreateCompanyBlock(ctx context.Context, input models.CompanyBlock) error

	GetCompanyBlockByID(ctx context.Context, ID uuid.UUID) (models.CompanyBlock, error)
	GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error)

	FilterCompanyBlocks(
		ctx context.Context,
		filters FilterParams,
		page, size uint64,
	) (models.CompanyBlockCollection, error)

	CancelActiveCompanyBlock(ctx context.Context, companyID uuid.UUID, canceledAt time.Time) error
}

func (s Service) getCompany(ctx context.Context, ID uuid.UUID) (models.Company, error) {
	company, err := s.db.GetCompanyByID(ctx, ID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by id %s, cause: %w", ID, err),
		)
	}

	if company.IsNil() {
		return models.Company{}, errx.ErrorcompanyNotFound.Raise(
			fmt.Errorf("company with ID %s not found", ID),
		)
	}

	return company, nil
}
