package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func TestDistributorBlocks(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := CreateDistributor(t, s)

	adminID := uuid.New()

	block, err := s.app.BlockDistributor(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("CreteBlock: %v", err)
	}

	if block.DistributorID != dis.ID {
		t.Errorf("expected blocked distributor ID '%s', got '%s'", dis.ID, block.DistributorID)
	}

	owner, err = s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.app.GetDistributor(ctx, dis.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dis.Status != enum.DistributorStatusBlocked {
		t.Errorf("expected distributor to be blocked")
	}
}

func TestUpdateBlockedDistributor(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	dis, owner := CreateDistributor(t, s)

	adminID := uuid.New()

	block, err := s.app.BlockDistributor(ctx, adminID, dis.ID, "Violation of terms")
	if err != nil {
		t.Fatalf("CreteBlock: %v", err)
	}

	if block.DistributorID != dis.ID {
		t.Errorf("expected blocked distributor ID '%s', got '%s'", dis.ID, block.DistributorID)
	}

	owner, err = s.app.GetEmployee(ctx, owner.UserID)
	if err != nil {
		t.Fatalf("GetEmployee: %v", err)
	}

	dis, err = s.app.GetDistributor(ctx, dis.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if dis.Status != enum.DistributorStatusBlocked {
		t.Errorf("expected distributor to be blocked")
	}

	name := "New Name"
	icon := "new_icon"
	_, err = s.app.UpdateDistributor(ctx, owner.UserID, dis.ID, app.UpdateDistributorParams{
		Name: &name,
		Icon: &icon,
	})
	if !errors.Is(err, errx.ErrorDistributorIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorDistributorIsBlocked, err)
	}

	_, err = s.app.SetDistributorStatus(ctx, owner.UserID, dis.ID, enum.DistributorStatusActive)
	if !errors.Is(err, errx.ErrorDistributorIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorDistributorIsBlocked, err)
	}

	_, err = s.app.SetDistributorStatus(ctx, owner.UserID, dis.ID, enum.DistributorStatusInactive)
	if !errors.Is(err, errx.ErrorDistributorIsBlocked) {
		t.Fatalf("expected error %v, got %v", errx.ErrorDistributorIsBlocked, err)
	}
}
