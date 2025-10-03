package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func TestCreateEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner := CreateDistributor(t, s)

	empOwn, err := s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	if empOwn.UserID != owner.UserID {
		t.Errorf("expected employee ID '%s', got '%s'", owner.UserID, empOwn.UserID)
	}

	if empOwn.Role != owner.Role {
		t.Errorf("expected employee role '%s', got '%s'", owner.Role, empOwn.Role)
	}

	emp, err := s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   owner.UserID,
		DistributorID: owner.DistributorID,
		Role:          enum.EmployeeRoleAdmin,
	})
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	if emp.Role != enum.EmployeeRoleAdmin {
		t.Errorf("expected employee role '%s', got '%s'", enum.EmployeeRoleAdmin, emp.Role)
	}
}

func TestErrorCreateEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner := CreateDistributor(t, s)

	empOwn, err := s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	if empOwn.UserID != owner.UserID {
		t.Errorf("expected employee ID '%s', got '%s'", owner.UserID, empOwn.UserID)
	}

	if empOwn.Role != owner.Role {
		t.Errorf("expected employee role '%s', got '%s'", owner.Role, empOwn.Role)
	}

	_, err = s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   owner.UserID,
		DistributorID: owner.DistributorID,
		Role:          enum.EmployeeRoleOwner,
	})
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorHaveNotEnoughRights, err)
	}

	_, err = s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   uuid.New(),
		DistributorID: owner.DistributorID,
		Role:          enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotEmployee) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotEmployee, err)
	}

	_, err = s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   owner.UserID,
		DistributorID: uuid.New(),
		Role:          enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotThisDistributorEmployee) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotThisDistributorEmployee, err)
	}
}

func TestInvalidDistributorInvite(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner1 := CreateDistributor(t, s)
	dist2, _ := CreateDistributor(t, s)

	_, err = s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   owner1.UserID,
		DistributorID: dist2.ID,
		Role:          enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotThisDistributorEmployee) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotThisDistributorEmployee, err)
	}
}

func TestCreateInvite(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	_, owner := CreateDistributor(t, s)

	inv, err := s.app.CreateInvite(ctx, app.CreateInviteParams{
		InitiatorID:   owner.UserID,
		DistributorID: owner.DistributorID,
		Role:          enum.EmployeeRoleAdmin,
	})
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}
	if inv.Status != enum.InviteStatusSent {
		t.Errorf("expected invite status '%s', got '%s'", enum.InviteStatusSent, inv.Status)
	}

	inv, err = s.app.AcceptInvite(ctx, uuid.New(), inv.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}
}
