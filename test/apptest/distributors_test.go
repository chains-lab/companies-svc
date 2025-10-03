package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func Createcompany(t *testing.T, s Setup) (models.Company, models.Employee) {
	ownerID := uuid.New()
	ctx := context.Background()

	dist, err := s.app.Createcompany(ctx, ownerID, "companyID 1", "icon")
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}

	owner, err := s.app.GetEmployee(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	return dist, owner
}

func CreateEmployee(t *testing.T, s Setup, initiatorID, companyID uuid.UUID, role string) models.Employee {
	ctx := context.Background()

	userID := uuid.New()

	imp, err := s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID: initiatorID,
		companyID:   companyID,
		Role:        role,
	})
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	inv, err := s.app.AcceptInvite(ctx, userID, imp.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	emp, err := s.app.GetEmployee(ctx, *inv.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	return emp
}

func TestCreatecompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	ownerID := uuid.New()
	dist, err := s.app.Createcompany(ctx, ownerID, "companyID 1", "icon")
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}

	dist, err = s.app.Getcompany(ctx, dist.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	owner, err := s.app.GetEmployee(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	owner, err = s.app.GetInitiator(ctx, ownerID)
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

func TestCreatecompanyByEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner := Createcompany(t, s)

	_, err = s.app.Createcompany(ctx, owner.UserID, "companyID 2", "icon2")
	if !errors.Is(err, errx.ErrorCurrentEmployeeCannotCreatecompany) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCurrentEmployeeCannotCreatecompany, err)
	}
}

func TestUpdatecompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := Createcompany(t, s)

	name := "Updated companyID Name"
	icon := "Updated Icon"
	dis, err = s.app.Updatecompany(ctx, owner.UserID, dis.ID, app.UpdatecompanyParams{
		Name: &name,
		Icon: &icon,
	})
	if err != nil {
		t.Fatalf("UpdateCompany: %v", err)
	}

	if dis.Name != name {
		t.Errorf("expected updated company name '%s', got '%s'", name, dis.Name)
	}
	if dis.Icon != icon {
		t.Errorf("expected updated company icon '%s', got '%s'", icon, dis.Icon)
	}
}

func TestInactivecompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := Createcompany(t, s)
	admin := CreateEmployee(t, s, owner.UserID, dis.ID, enum.EmployeeRoleAdmin)
	moder1 := CreateEmployee(t, s, admin.UserID, dis.ID, enum.EmployeeRoleModerator)
	moder2 := CreateEmployee(t, s, admin.UserID, dis.ID, enum.EmployeeRoleModerator)

	list, pag, err := s.app.ListEmployees(ctx, app.FilterEmployeeList{
		companies: []uuid.UUID{dis.ID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(list) != 4 {
		t.Fatalf("expected 4 employees, got %d", len(list))
	}
	if pag.Total != 4 {
		t.Fatalf("expected total 4 employees, got %d", pag.Total)
	}

	disID := dis.ID
	dis, err = s.app.Setcompaniestatus(ctx, owner.UserID, disID, enum.companiestatusBlocked)
	if !errors.Is(err, errx.ErrorUnexpectedcompaniesetStatus) {
		t.Fatalf("expected error %v, got %v", errx.ErrorUnexpectedcompaniesetStatus, err)
	}

	dis, err = s.app.Setcompaniestatus(ctx, admin.UserID, disID, enum.companiestatusBlocked)
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorHaveNotEnoughRights, err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, admin.UserID, disID)
	if !errors.Is(err, errx.ErrorCannotDeleteYourself) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, owner.UserID, disID)
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, moder1.UserID, moder2.UserID, disID)
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, moder2.UserID, disID)
	if err != nil {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	list, pag, err = s.app.ListEmployees(ctx, app.FilterEmployeeList{
		companies: []uuid.UUID{disID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 employee, got %d", len(list))
	}

	dis, err = s.app.Setcompaniestatus(ctx, owner.UserID, disID, enum.companiestatusInactive)
	if err != nil {
		t.Fatalf("Setcompaniestatus: %v", err)
	}
	if dis.Status != enum.companiestatusInactive {
		t.Errorf("expected company status '%s', got '%s'", enum.companiestatusInactive, dis.Status)
	}

	list, pag, err = s.app.ListEmployees(ctx, app.FilterEmployeeList{
		companies: []uuid.UUID{disID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 0 employee, got %d", len(list))
	}
	if pag.Total != 1 {
		t.Fatalf("expected total 1 employee, got %d", pag.Total)
	}
	if list[0].UserID != owner.UserID {
		t.Errorf("expected owner user id '%s', got '%s'", owner.UserID, list[0].UserID)
	}
}
