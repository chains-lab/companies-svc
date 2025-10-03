package company

import (
	"context"
	"time"

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

	CreateEmployee(ctx context.Context, input models.Employee) error
	GetUserEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error)

	CreateCompany(ctx context.Context, input models.Company) (models.Company, error)
	GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error)
	FilterCompanies(ctx context.Context, filters Filters, page, size uint64) (models.CompanyCollection, error)
	UpdateCompany(ctx context.Context, ID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	UpdateCompaniesStatus(ctx context.Context, ID uuid.UUID, status string, updatedAt time.Time) error

	CreateCompanyBlock(ctx context.Context, input models.CompanyBlock) error
	GetCompanyBlockByID(ctx context.Context, ID uuid.UUID) (models.CompanyBlock, error)
	GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error)
	FilterCompanyBlocks(ctx context.Context, filters FilterBlockages, page, size uint64) (models.CompanyBlockCollection, error)
	CancelActiveCompanyBlock(ctx context.Context, companyID uuid.UUID, canceledAt time.Time) error
}
