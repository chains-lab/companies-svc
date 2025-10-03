package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/internal/domain/service/invite"
	"github.com/chains-lab/companies-svc/test"

	"github.com/google/uuid"
)

func CreateCompany(t *testing.T, s Setup) (models.Company, models.Employee) {
	ownerID := uuid.New()
	ctx := context.Background()

	dist, err := s.domain.company.Create(ctx, ownerID, company.CreateParams{
		Name: "companyID 1",
		Icon: "icon",
	})
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}

	owner, err := s.domain.employee.Get(ctx, employee.GetParams{UserID: &ownerID})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	return dist, owner
}

func CreateEmployee(t *testing.T, s Setup, initiatorID, companyID uuid.UUID, role string) models.Employee {
	ctx := context.Background()

	userID := uuid.New()

	inv, err := s.domain.invite.Create(ctx, initiatorID, invite.CreateParams{
		CompanyID: companyID,
		Role:      role,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	inv, err = s.domain.invite.Accept(ctx, userID, inv.Token)
	if err != nil {
		t.Fatalf("Accept: %v", err)
	}

	emp, err := s.domain.employee.Get(ctx, employee.GetParams{UserID: inv.UserID})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	return emp
}

func TestCreateCompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	ownerID := uuid.New()
	dist, err := s.domain.company.Create(ctx, ownerID, company.CreateParams{
		Name: "companyID 1",
		Icon: "icon",
	})
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}

	dist, err = s.domain.company.Get(ctx, dist.ID)
	if err != nil {
		t.Fatalf("getCompany: %v", err)
	}

	owner, err := s.domain.employee.Get(ctx, employee.GetParams{UserID: &ownerID})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	owner, err = s.domain.employee.Get(ctx, employee.GetParams{UserID: &ownerID})
	if err != nil {
		t.Fatalf("GetInitiator: %v", err)
	}

	if owner.UserID != ownerID {
		t.Errorf("expected owner ID '%s', got '%s'", ownerID, owner.UserID)
	}

	if dist.Name != "companyID 1" {
		t.Errorf("expected company name 'companyID 1', got '%s'", dist.Name)
	}

	if owner.Role != enum.EmployeeRoleOwner {
		t.Errorf("expected owner role 'owner', got '%s'", owner.Role)
	}
}

func TestCreateCompanyByEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	_, owner := CreateCompany(t, s)

	_, err = s.domain.company.Create(ctx, owner.UserID, company.CreateParams{
		Name: "companyID 2",
		Icon: "icon 2",
	})
	if !errors.Is(err, errx.ErrorCurrentEmployeeCannotCreatecompany) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCurrentEmployeeCannotCreatecompany, err)
	}
}

func TestUpdateCompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	cou, _ := CreateCompany(t, s)

	name := "Updated companyID Name"
	icon := "Updated Icon"
	cou, err = s.domain.company.Update(ctx, cou.ID, company.UpdateParams{
		Name: &name,
		Icon: &icon,
	})
	if err != nil {
		t.Fatalf("UpdateCompany: %v", err)
	}

	if cou.Name != name {
		t.Errorf("expected updated company name '%s', got '%s'", name, cou.Name)
	}
	if cou.Icon != icon {
		t.Errorf("expected updated company icon '%s', got '%s'", icon, cou.Icon)
	}
}

func TestInactiveCompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	comp, owner := CreateCompany(t, s)
	admin := CreateEmployee(t, s, owner.UserID, comp.ID, enum.EmployeeRoleAdmin)
	moder1 := CreateEmployee(t, s, admin.UserID, comp.ID, enum.EmployeeRoleModerator)
	moder2 := CreateEmployee(t, s, admin.UserID, comp.ID, enum.EmployeeRoleModerator)

	moder1, err = s.domain.employee.Get(ctx, employee.GetParams{UserID: &moder1.UserID})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}
	if moder1.Role != enum.EmployeeRoleModerator {
		t.Errorf("expected moderator role '%s', got '%s'", enum.EmployeeRoleModerator, moder1.Role)
	}

	moder2, err = s.domain.employee.Get(ctx, employee.GetParams{UserID: &moder2.UserID})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}
	if moder2.Role != enum.EmployeeRoleModerator {
		t.Errorf("expected moderator role '%s', got '%s'", enum.EmployeeRoleModerator, moder2.Role)
	}

	emps, err := s.domain.employee.Filter(ctx, employee.FilterParams{
		CompanyID: &comp.ID,
	}, 1, 10)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(emps.Data) != 4 {
		t.Fatalf("expected 4 employees, got %d", len(emps.Data))
	}
	if emps.Total != 4 {
		t.Fatalf("expected total 4 employees, got %d", emps.Total)
	}

	compID := comp.ID
	comp, err = s.domain.company.UpdateStatus(ctx, compID, enum.DistributorStatusBlocked)
	if !errors.Is(err, errx.ErrorCannotSetcompaniestatusBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCannotSetcompaniestatusBlocked, err)
	}

	err = s.domain.employee.Delete(ctx, admin.UserID, admin.UserID, compID)
	if !errors.Is(err, errx.ErrorCannotDeleteYourself) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.domain.employee.Delete(ctx, admin.UserID, owner.UserID, compID)
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.domain.employee.Delete(ctx, moder1.UserID, moder2.UserID, compID)
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.domain.employee.Delete(ctx, admin.UserID, moder2.UserID, compID)
	if err != nil {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	emps, err = s.domain.employee.Filter(ctx, employee.FilterParams{
		CompanyID: &compID,
	}, 1, 10)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(emps.Data) != 3 {
		t.Fatalf("expected 3 employee, got %d", len(emps.Data))
	}

	comp, err = s.domain.company.UpdateStatus(ctx, compID, enum.DistributorStatusInactive)
	if err != nil {
		t.Fatalf("Setcompaniestatus: %v", err)
	}
	if comp.Status != enum.DistributorStatusInactive {
		t.Errorf("expected company status '%s', got '%s'", enum.DistributorStatusInactive, comp.Status)
	}

	err = s.domain.employee.RefuseOwn(ctx, admin.UserID)
	if err != nil {
		t.Fatalf("RefuseOwn: %v", err)
	}

	emps, err = s.domain.employee.Filter(ctx, employee.FilterParams{
		CompanyID: &compID,
	}, 1, 10)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(emps.Data) != 2 {
		t.Fatalf("expected 2 employee, got %d", len(emps.Data))
	}
	if emps.Total != 2 {
		t.Fatalf("expected total 2 employee, got %d", emps.Total)
	}
	if emps.Data[0].UserID != owner.UserID {
		t.Errorf("expected owner user id '%s', got '%s'", owner.UserID, emps.Data[0].UserID)
	}
}
