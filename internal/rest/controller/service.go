package controller

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/distributors-svc/internal/domain/service/employee"
	"github.com/chains-lab/logium"
	"github.com/google/uuid"
)

type DistributorSvc interface {
	CreteBlock(
		ctx context.Context,
		initiatorID uuid.UUID,
		distributorID uuid.UUID,
		reason string,
	) (models.DistributorBlock, error)
	FilterBlockages(
		ctx context.Context,
		filters distributor.FilterBlockages,
		page, size uint64,
	) (models.DistributorBlockCollection, error)
	GetBlock(ctx context.Context, BlockID uuid.UUID) (models.DistributorBlock, error)
	CancelBlock(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error)

	Create(ctx context.Context, params distributor.CreateParams) (models.Distributor, error)

	Filter(
		ctx context.Context,
		filters distributor.Filters,
		page, size uint64,
	) (models.DistributorCollection, error)

	Get(ctx context.Context, ID uuid.UUID) (models.Distributor, error)
	GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.DistributorBlock, error)

	UpdateStatus(
		ctx context.Context,
		distributorID uuid.UUID,
		status string,
	) (models.Distributor, error)

	Update(ctx context.Context,
		distributorID uuid.UUID,
		params distributor.UpdateParams,
	) (models.Distributor, error)
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

	Delete(ctx context.Context, initiatorID, userID, distributorID uuid.UUID) error
	RefuseOwn(ctx context.Context, initiatorID uuid.UUID) error
}

type domain struct {
	employee    EmployeeSvc
	distributor DistributorSvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, dis DistributorSvc, emp EmployeeSvc) Service {
	return Service{
		log: log,
		domain: domain{
			employee:    emp,
			distributor: dis,
		},
	}
}
