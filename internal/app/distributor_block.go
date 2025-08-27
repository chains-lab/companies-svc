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
	"github.com/chains-lab/pagi"
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
	FilterDistributorID(distributorID ...uuid.UUID) dbx.BlockQ
	FilterInitiatorID(initiatorID ...uuid.UUID) dbx.BlockQ
	FilterStatus(status ...string) dbx.BlockQ

	OrderByBlockedAt(ascending bool) dbx.BlockQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.BlockQ
}

func (a App) BlockDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	reason string,
) (models.DistributerBlock, error) {

	blockages := dbx.Blockages{
		ID:            uuid.New(),
		DistributorID: distributorID,
		InitiatorID:   initiatorID,
		Reason:        reason,
		Status:        enum.BlockStatusActive,
		BlockedAt:     time.Now().UTC(),
	}

	_, err := a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.DistributerBlock{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if err == nil {
		return models.DistributerBlock{}, errx.DistributorHaveAlreadyActiveBlock.Raise(
			fmt.Errorf("distributor %s already has an active block", distributorID),
		)
	}

	trErr := a.distributor.Transaction(func(ctx context.Context) error {
		err = a.distributor.New().FilterID(distributorID).Update(ctx, map[string]any{
			"status":     enum.DistributorStatusBlocked,
			"updated_at": time.Now().UTC(),
		})
		if err != nil {
			return errx.Internal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		_, err = a.block.New().FilterDistributorID(distributorID).FilterStatus(enum.BlockStatusActive).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errx.Internal.Raise(
				fmt.Errorf("checking existing block: %w", err),
			)
		}
		if err == nil {
			return errx.DistributorHaveAlreadyActiveBlock.Raise(
				fmt.Errorf("distributor %s already has an active block", distributorID),
			)
		}

		err = a.block.Insert(ctx, blockages)
		if err != nil {
			return errx.Internal.Raise(fmt.Errorf("inserting new block: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.DistributerBlock{}, trErr
	}

	return models.DistributerBlock{
		ID:            blockages.ID,
		DistributorID: blockages.DistributorID,
		InitiatorID:   blockages.InitiatorID,
		Reason:        blockages.Reason,
		BlockedAt:     blockages.BlockedAt,
	}, nil
}

func (a App) CanceledDistributorBlock(
	ctx context.Context,
	blockID uuid.UUID,
) (models.DistributerBlock, error) {
	block, err := a.block.New().FilterID(blockID).FilterStatus(enum.BlockStatusActive).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.DistributerBlock{}, errx.DistributorBlockNotFound.Raise(
				fmt.Errorf("block with ID %s not found", blockID),
			)
		default:
			return models.DistributerBlock{}, errx.Internal.Raise(
				fmt.Errorf("getting block with ID %s: %w", blockID, err),
			)
		}
	}
	canceledAt := time.Now().UTC()

	trErr := a.distributor.Transaction(func(ctx context.Context) error {
		err := a.block.FilterID(blockID).FilterStatus(enum.BlockStatusActive).Update(ctx, map[string]any{
			"status":       enum.BlockStatusCanceled,
			"cancelled_at": canceledAt,
		})
		if err != nil {
			return errx.Internal.Raise(
				fmt.Errorf("updating block status: %w", err),
			)
		}

		err = a.distributor.New().FilterID(block.DistributorID).Update(ctx, map[string]any{
			"status":     enum.DistributorStatusInactive,
			"updated_at": canceledAt,
		})
		if err != nil {
			return errx.Internal.Raise(
				fmt.Errorf("updating distributor status: %w", err),
			)
		}

		return nil
	})
	if trErr != nil {
		return models.DistributerBlock{}, trErr
	}

	res, err := a.block.New().FilterID(block.ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.DistributerBlock{},
				errx.DistributorBlockNotFound.Raise(
					fmt.Errorf("block with ID %s not found after unblock: %w", blockID, err),
				)
		default:
			return models.DistributerBlock{},
				errx.Internal.Raise(
					fmt.Errorf("getting block with ID %s after unblock: %w", blockID, err),
				)
		}
	}

	return models.DistributerBlock{
		ID:            res.ID,
		DistributorID: res.DistributorID,
		InitiatorID:   res.InitiatorID,
		Reason:        res.Reason,
		Status:        res.Status,
		BlockedAt:     res.BlockedAt,
		CanceledAt:    res.CanceledAt,
	}, nil
}

func (a App) GetBlock(
	ctx context.Context,
	ID uuid.UUID,
) (models.DistributerBlock, error) {
	block, err := a.block.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.DistributerBlock{}, errx.DistributorBlockNotFound.Raise(
				fmt.Errorf("block with ID %s not found: %w", ID, err),
			)
		default:
			return models.DistributerBlock{}, errx.Internal.Raise(
				fmt.Errorf("getting block with ID %s: %w", ID, err),
			)
		}
	}

	return models.DistributerBlock{
		ID:            block.ID,
		DistributorID: block.DistributorID,
		InitiatorID:   block.InitiatorID,
		Reason:        block.Reason,
		Status:        block.Status,
		BlockedAt:     block.BlockedAt,
		CanceledAt:    block.CanceledAt,
	}, nil
}

type SelectBlockagesParams struct {
	Distributors []uuid.UUID
	Initiators   []uuid.UUID
	Statuses     []string
}

func (a App) SelectBlockages(
	ctx context.Context,
	filters SelectBlockagesParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.DistributerBlock, pagi.Response, error) {
	query := a.block.New()

	if filters.Distributors != nil {
		query = query.FilterDistributorID(filters.Distributors...)
	}
	if filters.Initiators != nil {
		query = query.FilterInitiatorID(filters.Initiators...)
	}
	if filters.Statuses != nil {
		query = query.FilterStatus(filters.Statuses...)
	}

	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "blocked_at":
			query = query.OrderByBlockedAt(ascend)
		default:

		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("counting rows: %w", err),
		)
	}

	rows, err := query.OrderByBlockedAt(true).Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("selecting rows: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	res := make([]models.DistributerBlock, 0, len(rows))
	for _, s := range rows {
		el := models.DistributerBlock{
			ID:            s.ID,
			DistributorID: s.DistributorID,
			InitiatorID:   s.InitiatorID,
			Reason:        s.Reason,
			Status:        s.Status,
			BlockedAt:     s.BlockedAt,
		}
		if s.CanceledAt != nil {
			canceledAt := *s.CanceledAt
			el.CanceledAt = &canceledAt
		}

		res = append(res, el)
	}

	return res, pagi.Response{
		Total: count,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}
