package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

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
