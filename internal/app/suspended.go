package app

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"github.com/google/uuid"
)

type suspendedQ interface {
	New() dbx.SuspendedQ

	Insert(ctx context.Context, input dbx.SuspendedDistributor) error
	Get(ctx context.Context) (dbx.SuspendedDistributor, error)
	Select(ctx context.Context) ([]dbx.SuspendedDistributor, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.SuspendedQ
	FilterDistributorID(distributorID uuid.UUID) dbx.SuspendedQ
	FilterInitiatorID(initiatorID uuid.UUID) dbx.SuspendedQ
	FilterActive(active bool) dbx.SuspendedQ

	OrderBySuspendedAt(ascending bool) dbx.SuspendedQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.SuspendedQ
}

func (a App) GetSuspendsDistributor(
	ctx context.Context,
	filters map[string]any,
	pag pagination.Request) (
	//returns:
	[]models.SuspendedDistributor,
	pagination.Response,
	error,
) {
	query := a.suspended.New()

	if id, ok := filters["id"].(uuid.UUID); ok {
		query = query.FilterID(id)
	}
	if distributorID, ok := filters["distributor_id"].(uuid.UUID); ok {
		query = query.FilterDistributorID(distributorID)
	}
	if initiatorID, ok := filters["initiator_id"].(uuid.UUID); ok {
		query = query.FilterInitiatorID(initiatorID)
	}
	if active, ok := filters["active"].(bool); ok {
		query = query.FilterActive(active)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	suspendeds, err := query.OrderBySuspendedAt(true).Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	res := make([]models.SuspendedDistributor, len(suspendeds))
	for _, s := range suspendeds {
		el := models.SuspendedDistributor{
			ID:            s.ID,
			DistributorID: s.DistributorID,
			InitiatorID:   s.InitiatorID,
			Reason:        s.Reason,
			Active:        s.Active,
			SuspendedAt:   s.SuspendedAt,
			CreatedAt:     s.CreatedAt,
		}
		if s.CanceledAt != nil {
			canceledAt := *s.CanceledAt
			el.CanceledAt = &canceledAt
		}

		res = append(res, el)
	}

	return res, pagination.Response{
		Total: count,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}

func (a App) CreateSuspendedDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.SuspendedDistributor, error) {
	if initiatorID == uuid.Nil || distributorID == uuid.Nil {
		return models.SuspendedDistributor{}, errx.RaiseBadRequest(ctx, "initiator_id and distributor_id cannot be empty")
	}

	suspended := dbx.SuspendedDistributor{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Active:        true,
		SuspendedAt:   time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	}

	err := a.suspended.Insert(ctx, suspended)
	if err != nil {
		return models.SuspendedDistributor{}, errx.RaiseInternal(ctx, err)
	}

	return models.SuspendedDistributor{
		ID:            suspended.ID,
		DistributorID: suspended.DistributorID,
		InitiatorID:   suspended.InitiatorID,
		Reason:        suspended.Reason,
		SuspendedAt:   suspended.SuspendedAt,
		CreatedAt:     suspended.CreatedAt,
	}, nil
}

func (a App) CancelSuspendedDistributor(
	ctx context.Context,
	id uuid.UUID,
) error {
	if id == uuid.Nil {
		return errx.RaiseBadRequest(ctx, "id cannot be empty")
	}

	query := a.suspended.New().FilterID(id)

	suspended, err := query.Get(ctx)
	if err != nil {
		return errx.RaiseInternal(ctx, err)
	}

	if !suspended.Active {
		return errx.RaiseBadRequest(ctx, "suspended distributor is already canceled")
	}

	err = query.Update(ctx, map[string]any{
		"active":      false,
		"canceled_at": time.Now().UTC(),
	})
	if err != nil {
		return errx.RaiseInternal(ctx, err)
	}

	return nil
}
