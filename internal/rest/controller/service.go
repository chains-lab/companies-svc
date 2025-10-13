package controller

import (
	"context"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/block"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/domain/service/invite"
	"github.com/chains-lab/logium"
	"github.com/google/uuid"
)

type companySvc interface {
	Create(ctx context.Context, initiatorID uuid.UUID, params company.CreateParams) (models.Company, error)

	Get(ctx context.Context, ID uuid.UUID) (models.Company, error)
	Filter(
		ctx context.Context,
		filters company.FiltersParams,
		page, size uint64,
	) (models.CompanyCollection, error)

	Update(ctx context.Context, companyID uuid.UUID, params company.UpdateParams) (models.Company, error)
	UpdateStatus(
		ctx context.Context,
		companyID uuid.UUID,
		status string,
	) (models.Company, error)
}

type employeeSvc interface {
	Create(ctx context.Context, params employee.CreateParams) (models.Employee, error)

	Get(ctx context.Context, params employee.GetParams) (models.Employee, error)
	Filter(
		ctx context.Context,
		filters employee.FilterParams,
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

type inviteSvc interface {
	Create(ctx context.Context, InitiatorID uuid.UUID, params invite.CreateParams) (models.Invite, error)
	Accept(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error)
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
	) (models.CompanyBlockCollection, error)

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
