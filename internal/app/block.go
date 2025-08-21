package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"github.com/google/uuid"
)

type blockagesQ interface {
	New() dbx.BlockQ

	Insert(ctx context.Context, input dbx.Blockages) error
	Get(ctx context.Context) (dbx.Blockages, error)
	Select(ctx context.Context) ([]dbx.Blockages, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.BlockQ
	FilterDistributorID(distributorID uuid.UUID) dbx.BlockQ
	FilterInitiatorID(initiatorID uuid.UUID) dbx.BlockQ
	FilterStatus(status string) dbx.BlockQ
	OrderByBlockedAt(ascending bool) dbx.BlockQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.BlockQ
}

func (a App) CreateDistributorBlock(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.Block, error) {

	blockages := dbx.Blockages{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Status:        enum.BlockStatusActive,
		BlockedAt:     time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	}

	_, err := a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Block{}, errx.RaiseInternal(ctx, fmt.Errorf("checking existing block: %w", err))
	}
	if err == nil {
		return models.Block{}, errx.RaiseDistributorHaveAlreadyActiveBlock(
			ctx,
			fmt.Errorf("distributor %s already has an active block", distributorID),
			distributorID,
		)
	}

	trErr := a.distributor.Transaction(func(ctx context.Context) error {
		err = a.distributor.New().FilterID(distributorID).Update(ctx, map[string]any{
			"status":     enum.DistributorStatusBlocked,
			"updated_at": time.Now().UTC(),
		})
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("updating distributor status: %w", err))
		}

		_, err = a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errx.RaiseInternal(ctx, fmt.Errorf("checking existing block: %w", err))
		}
		if err == nil {
			return errx.RaiseDistributorHaveAlreadyActiveBlock(
				ctx,
				fmt.Errorf("distributor %s already has an active block", distributorID),
				distributorID,
			)
		}

		err = a.block.Insert(ctx, blockages)
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("inserting new block: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.Block{}, trErr
	}

	return models.Block{
		ID:            blockages.ID,
		DistributorID: blockages.DistributorID,
		InitiatorID:   blockages.InitiatorID,
		Reason:        blockages.Reason,
		BlockedAt:     blockages.BlockedAt,
		CreatedAt:     blockages.CreatedAt,
	}, nil
}

func (a App) CancelBlockForDistributor(
	ctx context.Context,
	distributorID uuid.UUID,
) (models.Block, error) {
	block, err := a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Block{}, errx.RaiseDistributorHaveNotActiveBlock(
				ctx,
				fmt.Errorf("distributor %s have not active block", distributorID),
				distributorID,
			)
		default:
			return models.Block{},
				errx.RaiseInternal(ctx, fmt.Errorf("getting block for distributor %s: %w", distributorID, err))
		}
	}

	blockQ := a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive)

	canceledAt := time.Now().UTC()

	trErr := a.distributor.Transaction(func(ctx context.Context) error {
		err := blockQ.Update(ctx, map[string]any{
			"status":       enum.BlockStatusCanceled,
			"cancelled_at": canceledAt,
		})
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("updating block status: %w", err))
		}

		distributorQ := a.distributor.New().FilterID(distributorID)

		err = distributorQ.Update(ctx, map[string]any{
			"status":     enum.DistributorStatusInactive,
			"updated_at": canceledAt,
		})
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("updating distributor status: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.Block{}, trErr
	}

	res, err := a.block.New().FilterID(block.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Block{},
				errx.RaiseBlockNotFound(ctx, fmt.Errorf("block for distributor %s not found", distributorID), distributorID)
		default:
			return models.Block{},
				errx.RaiseInternal(ctx, fmt.Errorf("getting block for distributor %s: %w", distributorID, err))
		}
	}

	return models.Block{
		ID:            res.ID,
		DistributorID: res.DistributorID,
		InitiatorID:   res.InitiatorID,
		Reason:        res.Reason,
		Status:        res.Status,
		BlockedAt:     res.BlockedAt,
		CanceledAt:    res.CanceledAt,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (a App) GetBlock(
	ctx context.Context,
	ID uuid.UUID,
) (models.Block, error) {
	block, err := a.block.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Block{}, errx.RaiseBlockNotFound(
				ctx,
				fmt.Errorf("block with ID %s not found", ID),
				ID,
			)
		default:
			return models.Block{}, errx.RaiseInternal(ctx, fmt.Errorf("getting block with ID %s: %w", ID, err))
		}
	}

	return models.Block{
		ID:            block.ID,
		DistributorID: block.DistributorID,
		InitiatorID:   block.InitiatorID,
		Reason:        block.Reason,
		Status:        block.Status,
		BlockedAt:     block.BlockedAt,
		CanceledAt:    block.CanceledAt,
		CreatedAt:     block.CreatedAt,
	}, nil
}

func (a App) SelectBlockages(
	ctx context.Context,
	filters map[string]any,
	pag pagination.Request,
) ([]models.Block, pagination.Response, error) {
	query := a.block.New()

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

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, fmt.Errorf("counting blockages: %w", err))
	}

	blockages, err := query.OrderByBlockedAt(true).Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, fmt.Errorf("selecting blockages: %w", err))
	}

	res := make([]models.Block, 0, len(blockages))
	for _, s := range blockages {
		el := models.Block{
			ID:            s.ID,
			DistributorID: s.DistributorID,
			InitiatorID:   s.InitiatorID,
			Reason:        s.Reason,
			Status:        s.Status,
			BlockedAt:     s.BlockedAt,
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
