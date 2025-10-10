package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/companies-svc/internal/domain/service/employee"
	"github.com/chains-lab/companies-svc/test"
	"github.com/google/uuid"
)

func TestCompanyBlocks(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	dis, owner := CreateCompany(t, s)

	adminID := uuid.New()

	block, err := s.domain.block.Crete(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("Crete: %v", err)
	}
	if block.CompanyID != dis.ID {
		t.Errorf("expected blocked company ID '%s', got '%s'", dis.ID, block.CompanyID)
	}

	owner, err = s.domain.employee.Get(ctx, employee.GetParams{
		UserID: &owner.UserID,
	})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.domain.company.Get(ctx, dis.ID)
	if err != nil {
		t.Fatalf("getCompany: %v", err)
	}
	if dis.Status != enum.CompanyStatusBlocked {
		t.Errorf("expected company to be blocked")
	}
}

func TestUpdateBlockedCompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	dis, owner := CreateCompany(t, s)

	adminID := uuid.New()

	block, err := s.domain.block.Crete(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("Crete: %v", err)
	}

	if block.CompanyID != dis.ID {
		t.Errorf("expected blocked company ID '%s', got '%s'", dis.ID, block.CompanyID)
	}

	owner, err = s.domain.employee.Get(ctx, employee.GetParams{
		UserID: &owner.UserID,
	})
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.domain.company.Get(ctx, dis.ID)
	if err != nil {
		t.Fatalf("getCompany: %v", err)
	}
	if dis.Status != enum.CompanyStatusBlocked {
		t.Errorf("expected company to be blocked")
	}

	name := "New Name"
	icon := "new_icon"
	_, err = s.domain.company.Update(ctx, dis.ID, company.UpdateParams{
		Name: &name,
		Icon: &icon,
	})
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}

	_, err = s.domain.company.UpdateStatus(ctx, dis.ID, enum.CompanyStatusActive)
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}

	_, err = s.domain.company.UpdateStatus(ctx, dis.ID, enum.CompanyStatusInactive)
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}
}
