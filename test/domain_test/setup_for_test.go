package domain_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/chains-lab/companies-svc/internal"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/services/block"
	"github.com/chains-lab/companies-svc/internal/domain/services/company"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/domain/services/invite"
	"github.com/chains-lab/companies-svc/internal/jwtmanager"
	"github.com/chains-lab/companies-svc/internal/repo"
	"github.com/chains-lab/companies-svc/internal/usrguesser"
	"github.com/chains-lab/companies-svc/test"
	"github.com/google/uuid"
)

type companySvc interface {
	create(ctx context.Context, initiatorID uuid.UUID, params company.CreateParams) (models.Company, error)

	Get(ctx context.Context, ID uuid.UUID) (models.Company, error)
	Filter(
		ctx context.Context,
		filters company.FiltersParams,
		page, size uint64,
	) (models.CompaniesCollection, error)

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
	) (models.EmployeesCollection, error)

	UpdateEmployeeRole(
		ctx context.Context,
		initiatorID uuid.UUID,
		userID uuid.UUID,
		newRole string,
	) (models.Employee, error)

	Delete(ctx context.Context, initiatorID, userID, companyID uuid.UUID) error
	RefuseMe(ctx context.Context, initiatorID uuid.UUID) error
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
	) (models.CompanyBlocksCollection, error)

	Get(ctx context.Context, blockID uuid.UUID) (models.CompanyBlock, error)
	GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error)
}

type domain struct {
	company  companySvc
	employee employeeSvc
	invite   inviteSvc
	block    blockSvc
}
type Setup struct {
	domain domain
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := internal.Config{
		JWT: internal.JWTConfig{
			Invites: struct {
				SecretKey string `mapstructure:"secret_key"`
			}{
				SecretKey: "invitesuperkey", // тут подставь ключ для тестов
			},
		},
		Database: internal.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: test.TestDatabaseURL,
			},
		},
		Profile: internal.ProfileConfig{
			Url: "http://localhost:8001/profiles-svc/v1/profiles",
		},
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := repo.NewDatabase(pg)
	jwtInviteManager := jwtmanager.NewManager(cfg)
	userGuesser := usrguesser.NewService(cfg.Profile.Url, nil)

	companiesSvc := company.NewService(database)
	employeeSvc := employee.NewService(database, userGuesser)
	inviteSvc := invite.NewService(database, jwtInviteManager)
	blockSvc := block.NewService(database)

	return Setup{
		domain: domain{
			company:  companiesSvc,
			employee: employeeSvc,
			invite:   inviteSvc,
			block:    blockSvc,
		},
	}, nil
}
