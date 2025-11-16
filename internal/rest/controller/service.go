package controller

import (
	"context"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/services/block"
	"github.com/chains-lab/companies-svc/internal/domain/services/company"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/domain/services/invite"
	"github.com/chains-lab/logium"
	"github.com/google/uuid"
)

type companySvc interface {
	CreateByEmployee(
		ctx context.Context,
		initiatorID uuid.UUID,
		params company.CreateParams,
	) (models.Company, error)

	Get(ctx context.Context, ID uuid.UUID) (models.Company, error)
	Filter(
		ctx context.Context,
		filters company.FiltersParams,
		page, size uint64,
	) (models.CompaniesCollection, error)

	UpdateByEmployee(
		ctx context.Context,
		initiatorID, companyID uuid.UUID,
		params company.UpdateParams,
	) (models.Company, error)
	UpdateStatusByEmployee(
		ctx context.Context,
		initiatorID, companyID uuid.UUID,
		status string,
	) (models.Company, error)
}

type employeeSvc interface {
	Get(ctx context.Context, params employee.GetParams) (models.Employee, error)
	GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.Employee, error)

	Filter(
		ctx context.Context,
		filters employee.FilterParams,
		page, size uint64,
	) (models.EmployeesCollection, error)

	UpdateByEmployee(
		ctx context.Context,
		userID uuid.UUID,
		initiatorID uuid.UUID,
		params employee.UpdateParams,
	) (models.Employee, error)
	UpdateMy(
		ctx context.Context,
		initiatorID uuid.UUID,
		params employee.UpdateMyParams,
	) (models.Employee, error)

	DeleteByEmployee(
		ctx context.Context,
		initiatorID, userID, companyID uuid.UUID,
	) error
	DeleteMe(
		ctx context.Context,
		initiatorID uuid.UUID,
	) error
}

type inviteSvc interface {
	Create(ctx context.Context, initiatorID uuid.UUID, params invite.CreateParams) (models.Invite, error)
	Answer(ctx context.Context, userID, inviteID uuid.UUID, answer string) (models.Invite, error)
}

type blockSvc interface {
	Crete(
		ctx context.Context,
		initiatorID uuid.UUID,
		companyID uuid.UUID,
		reason string,
	) (models.CompanyBlock, error)

	Cancel(ctx context.Context, companyID uuid.UUID) (models.Company, error)

	Filter(
		ctx context.Context,
		filters block.FilterParams,
		page, size uint64,
	) (models.CompanyBlocksCollection, error)

	Get(ctx context.Context, blockID uuid.UUID) (models.CompanyBlock, error)
	GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error)
}

type domain struct {
	employee employeeSvc
	company  companySvc
	invite   inviteSvc
	block    blockSvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, comp companySvc, emp employeeSvc, inv inviteSvc, blc blockSvc) Service {
	return Service{
		log: log,
		domain: domain{
			employee: emp,
			company:  comp,
			invite:   inv,
			block:    blc,
		},
	}
}
