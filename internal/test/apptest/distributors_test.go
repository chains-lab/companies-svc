package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func CreateDistributor(t *testing.T, s Setup) (models.Distributor, models.Employee) {
	ownerID := uuid.New()
	ctx := context.Background()

	dist, err := s.app.CreateDistributor(ctx, ownerID, "Distributor 1", "icon")
	if err != nil {
		t.Fatalf("CreateDistributor: %v", err)
	}

	owner, err := s.app.GetEmployee(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	return dist, owner
}

func CreateEmployee(t *testing.T, s Setup, initiatorID, distributorID uuid.UUID, role string) models.Employee {
	ctx := context.Background()

	userID := uuid.New()

	imp, err := s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   initiatorID,
		DistributorID: distributorID,
		Role:          role,
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

func TestCreateDistributor(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	ownerID := uuid.New()
	dist, err := s.app.CreateDistributor(ctx, ownerID, "Distributor 1", "icon")
	if err != nil {
		t.Fatalf("CreateDistributor: %v", err)
	}

	dist, err = s.app.GetDistributor(ctx, dist.ID)
	if err != nil {
		t.Fatalf("GetDistributor: %v", err)
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

	if dist.Name != "Distributor 1" {
		t.Errorf("expected distributor name 'Distributor 1', got '%s'", dist.Name)
	}

	if owner.Role != enum.EmployeeRoleOwner {
		t.Errorf("expected owner role 'owner', got '%s'", owner.Role)
	}
}

func TestCreateDistributorByEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner := CreateDistributor(t, s)

	_, err = s.app.CreateDistributor(ctx, owner.UserID, "Distributor 2", "icon2")
	if !errors.Is(err, errx.ErrorCurrentEmployeeCannotCreateDistributor) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCurrentEmployeeCannotCreateDistributor, err)
	}
}

func TestUpdateDistributor(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := CreateDistributor(t, s)

	name := "Updated Distributor Name"
	icon := "Updated Icon"
	dis, err = s.app.UpdateDistributor(ctx, owner.UserID, dis.ID, app.UpdateDistributorParams{
		Name: &name,
		Icon: &icon,
	})
	if err != nil {
		t.Fatalf("UpdateDistributor: %v", err)
	}

	if dis.Name != name {
		t.Errorf("expected updated distributor name '%s', got '%s'", name, dis.Name)
	}
	if dis.Icon != icon {
		t.Errorf("expected updated distributor icon '%s', got '%s'", icon, dis.Icon)
	}
}

func TestInactiveDistributor(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := CreateDistributor(t, s)
	admin := CreateEmployee(t, s, owner.UserID, dis.ID, enum.EmployeeRoleAdmin)
	moder1 := CreateEmployee(t, s, admin.UserID, dis.ID, enum.EmployeeRoleModerator)
	moder2 := CreateEmployee(t, s, admin.UserID, dis.ID, enum.EmployeeRoleModerator)

	list, pag, err := s.app.ListEmployees(ctx, app.FilterEmployeeList{
		Distributors: []uuid.UUID{dis.ID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("ListEmployees: %v", err)
	}
	if len(list) != 4 {
		t.Fatalf("expected 4 employees, got %d", len(list))
	}
	if pag.Total != 4 {
		t.Fatalf("expected total 4 employees, got %d", pag.Total)
	}

	disID := dis.ID
	dis, err = s.app.SetDistributorStatus(ctx, owner.UserID, disID, enum.DistributorStatusBlocked)
	if !errors.Is(err, errx.ErrorUnexpectedDistributorSetStatus) {
		t.Fatalf("expected error %v, got %v", errx.ErrorUnexpectedDistributorSetStatus, err)
	}

	dis, err = s.app.SetDistributorStatus(ctx, admin.UserID, disID, enum.DistributorStatusBlocked)
	if !errors.Is(err, errx.ErrorInitiatorEmployeeHaveNotEnoughRights) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorEmployeeHaveNotEnoughRights, err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, admin.UserID, disID)
	if !errors.Is(err, errx.ErrorCannotDeleteYourself) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, owner.UserID, disID)
	if !errors.Is(err, errx.ErrorInitiatorEmployeeHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, moder1.UserID, moder2.UserID, disID)
	if !errors.Is(err, errx.ErrorInitiatorEmployeeHaveNotEnoughRights) {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	err = s.app.DeleteEmployee(ctx, admin.UserID, moder2.UserID, disID)
	if err != nil {
		t.Fatalf("DeleteEmployee: %v", err)
	}

	list, pag, err = s.app.ListEmployees(ctx, app.FilterEmployeeList{
		Distributors: []uuid.UUID{disID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("ListEmployees: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 employee, got %d", len(list))
	}

	dis, err = s.app.SetDistributorStatus(ctx, owner.UserID, disID, enum.DistributorStatusInactive)
	if err != nil {
		t.Fatalf("SetDistributorStatus: %v", err)
	}
	if dis.Status != enum.DistributorStatusInactive {
		t.Errorf("expected distributor status '%s', got '%s'", enum.DistributorStatusInactive, dis.Status)
	}

	list, pag, err = s.app.ListEmployees(ctx, app.FilterEmployeeList{
		Distributors: []uuid.UUID{disID},
	}, pagi.Request{
		Page: 1,
		Size: 10,
	}, nil)
	if err != nil {
		t.Fatalf("ListEmployees: %v", err)
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
