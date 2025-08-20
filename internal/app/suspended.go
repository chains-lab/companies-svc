package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
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
	FilterStatus(status string) dbx.SuspendedQ
	OrderBySuspendedAt(ascending bool) dbx.SuspendedQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.SuspendedQ
}

func (a App) GetSuspendedDistributor(
	ctx context.Context,
	id uuid.UUID,
) (models.SuspendedDistributor, error) {
	suspended, err := a.suspended.New().FilterID(id).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.SuspendedDistributor{}, errx.RaiseSuspendedDistributorNotFound(ctx, id)
		default:
			return models.SuspendedDistributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.SuspendedDistributor{
		ID:            suspended.ID,
		DistributorID: suspended.DistributorID,
		InitiatorID:   suspended.InitiatorID,
		Reason:        suspended.Reason,
		Status:        suspended.Status,
		SuspendedAt:   suspended.SuspendedAt,
		CanceledAt:    suspended.CanceledAt,
		CreatedAt:     suspended.CreatedAt,
	}, nil
}

func (a App) GetSuspendsDistributor(
	ctx context.Context,
	filters map[string]any,
	pag pagination.Request,
) ([]models.SuspendedDistributor, pagination.Response, error) {
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
	if status, ok := filters["status"].(string); ok {
		query = query.FilterStatus(status)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	suspendeds, err := query.OrderBySuspendedAt(true).Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, nil // No suspended distributors found
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, nil // No suspended distributors found
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	res := make([]models.SuspendedDistributor, len(suspendeds))
	for _, s := range suspendeds {
		el := models.SuspendedDistributor{
			ID:            s.ID,
			DistributorID: s.DistributorID,
			InitiatorID:   s.InitiatorID,
			Reason:        s.Reason,
			Status:        s.Status,
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
	suspended := dbx.SuspendedDistributor{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Status:        enum.SuspendedDistributorStatusActive,
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
	query := a.suspended.New().FilterID(id)

	err := query.Update(ctx, map[string]any{
		"status":      enum.SuspendedDistributorStatusCanceled,
		"canceled_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseSuspendedDistributorNotFound(ctx, id)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	return nil
}
