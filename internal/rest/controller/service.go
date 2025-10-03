package controller

import (
	"context"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/logium"
	"github.com/google/uuid"
)

type companiesvc interface {
	CreteBlock(
		ctx context.Context,
		initiatorID uuid.UUID,
		companyID uuid.UUID,
		reason string,
	) (models.CompanyBlock, error)
	FilterBlockages(
		ctx context.Context,
		filters company.FilterBlockages,
		page, size uint64,
	) (models.CompanyBlockCollection, error)
	GetBlock(ctx context.Context, BlockID uuid.UUID) (models.CompanyBlock, error)
	CancelBlock(ctx context.Context, companyID uuid.UUID) (models.Company, error)

	Create(ctx context.Context, params company.CreateParams) (models.Company, error)

	Filter(
		ctx context.Context,
		filters company.Filters,
		page, size uint64,
	) (models.CompanyCollection, error)

	Get(ctx context.Context, ID uuid.UUID) (models.Company, error)
	GetActivecompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error)

	UpdateStatus(
		ctx context.Context,
		companyID uuid.UUID,
		status string,
	) (models.Company, error)

	Update(ctx context.Context,
		companyID uuid.UUID,
		params company.UpdateParams,
	) (models.Company, error)
}

type EmployeeSvc interface {
	CreateInvite(ctx context.Context, InitiatorID uuid.UUID, params employee.SentInviteParams) (models.Invite, error)
	AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error)
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)

	Create(ctx context.Context, params employee.CreateParams) (models.Employee, error)

	Get(ctx context.Context, filters employee.GetFilters) (models.Employee, error)
	Filter(
		ctx context.Context,
		filters employee.Filter,
		page, size uint64,
	) (models.EmployeeCollection, error)

	UpdateEmployeeRole(
		ctx context.Context,
		initiatorID uuid.UUID,
		userID uuid.UUID,
		newRole string,
	) (models.Employee, error)

	Delete(ctx context.Context, initiatorID, userID, companyID uuid.UUID) error
	RefuseOwn(ctx context.Context, initiatorID uuid.UUID) error
}

type domain struct {
	employee EmployeeSvc
	company  companiesvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, dis companiesvc, emp EmployeeSvc) Service {
	return Service{
		log: log,
		domain: domain{
			employee: emp,
			company:  dis,
		},
	}
}
