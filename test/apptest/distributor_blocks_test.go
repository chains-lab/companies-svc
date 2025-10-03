package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func TestcompanyBlocks(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := Createcompany(t, s)

	adminID := uuid.New()

	block, err := s.app.Blockcompany(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("CreteBlock: %v", err)
	}

	if block.CompanyID != dis.ID {
		t.Errorf("expected blocked company ID '%s', got '%s'", dis.ID, block.CompanyID)
	}

	owner, err = s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.app.Getcompany(ctx, dis.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dis.Status != enum.companiestatusBlocked {
		t.Errorf("expected company to be blocked")
	}
}

func TestUpdateBlockedcompany(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := Createcompany(t, s)

	adminID := uuid.New()

	block, err := s.app.Blockcompany(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("CreteBlock: %v", err)
	}

	if block.CompanyID != dis.ID {
		t.Errorf("expected blocked company ID '%s', got '%s'", dis.ID, block.CompanyID)
	}

	owner, err = s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.app.Getcompany(ctx, dis.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dis.Status != enum.companiestatusBlocked {
		t.Errorf("expected company to be blocked")
	}

	name := "New Name"
	icon := "new_icon"
	_, err = s.app.Updatecompany(ctx, owner.UserID, dis.ID, app.UpdatecompanyParams{
		Name: &name,
		Icon: &icon,
	})
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}

	_, err = s.app.Setcompaniestatus(ctx, owner.UserID, dis.ID, enum.companiestatusActive)
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}

	_, err = s.app.Setcompaniestatus(ctx, owner.UserID, dis.ID, enum.companiestatusInactive)
	if !errors.Is(err, errx.ErrorcompanyIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorcompanyIsBlocked, err)
	}
}
