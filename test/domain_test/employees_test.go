package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/services/employee"
	"github.com/chains-lab/companies-svc/internal/domain/services/invite"
	"github.com/chains-lab/companies-svc/test"
	"github.com/google/uuid"
)

func TestCreateEmployee(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	_, owner := CreateCompany(t, s)

	empOwn, err := s.domain.employee.Get(ctx, employee.GetParams{
		UserID: &owner.UserID,
	})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	if empOwn.UserID != owner.UserID {
		t.Errorf("expected employee ID '%s', got '%s'", owner.UserID, empOwn.UserID)
	}

	if empOwn.Role != owner.Role {
		t.Errorf("expected employee role '%s', got '%s'", owner.Role, empOwn.Role)
	}

	emp, err := s.domain.invite.Create(ctx, owner.UserID, invite.CreateParams{
		CompanyID: owner.CompanyID,
		Role:      enum.EmployeeRoleAdmin,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
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

	test.CleanDb(t)

	ctx := context.Background()

	_, owner := CreateCompany(t, s)

	empOwn, err := s.domain.employee.Get(ctx, employee.GetParams{
		UserID: &owner.UserID,
	})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	if empOwn.UserID != owner.UserID {
		t.Errorf("expected employee ID '%s', got '%s'", owner.UserID, empOwn.UserID)
	}

	if empOwn.Role != owner.Role {
		t.Errorf("expected employee role '%s', got '%s'", owner.Role, empOwn.Role)
	}

	_, err = s.domain.invite.Create(ctx, owner.UserID, invite.CreateParams{
		CompanyID: owner.CompanyID,
		Role:      enum.EmployeeRoleOwner,
	})
	if !errors.Is(err, errx.ErrorInitiatorHaveNotEnoughRights) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorHaveNotEnoughRights, err)
	}

	_, err = s.domain.invite.Create(ctx, uuid.New(), invite.CreateParams{
		CompanyID: owner.CompanyID,
		Role:      enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotEmployee) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotEmployee, err)
	}

	_, err = s.domain.invite.Create(ctx, owner.UserID, invite.CreateParams{
		CompanyID: uuid.New(),
		Role:      enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotEmployeeOfThisCompany) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotEmployeeOfThisCompany, err)
	}
}

func TestInvalidCompanyInvite(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	_, owner1 := CreateCompany(t, s)
	comp2, _ := CreateCompany(t, s)

	_, err = s.domain.invite.Create(ctx, owner1.UserID, invite.CreateParams{
		CompanyID: comp2.ID,
		Role:      enum.EmployeeRoleAdmin,
	})
	if !errors.Is(err, errx.ErrorInitiatorIsNotEmployeeOfThisCompany) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInitiatorIsNotEmployeeOfThisCompany, err)
	}
}

func TestCreateInvite(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	_, owner := CreateCompany(t, s)

	inv, err := s.domain.invite.Create(ctx, owner.UserID, invite.CreateParams{
		CompanyID: owner.CompanyID,
		Role:      enum.EmployeeRoleAdmin,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if inv.Status != enum.InviteStatusSent {
		t.Errorf("expected invite status '%s', got '%s'", enum.InviteStatusSent, inv.Status)
	}

	inv, err = s.domain.invite.Accept(ctx, uuid.New(), inv.Token)
	if err != nil {
		t.Fatalf("Accept: %v", err)
	}
}
