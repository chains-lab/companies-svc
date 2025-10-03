package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const blockedTable = "distributor_blockages"

type DistributorBlock struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID uuid.UUID  `db:"distributor_id"`
	InitiatorID   uuid.UUID  `db:"initiator_id"`
	Reason        string     `db:"reason"`
	Status        string     `db:"status"`
	BlockedAt     time.Time  `db:"blocked_at"`
	CanceledAt    *time.Time `db:"canceled_at"`
}

type BlockagesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewBlockagesQ(db *sql.DB) BlockagesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return BlockagesQ{
		db:       db,
		selector: builder.Select("*").From(blockedTable),
		updater:  builder.Update(blockedTable),
		inserter: builder.Insert(blockedTable),
		deleter:  builder.Delete(blockedTable),
		counter:  builder.Select("COUNT(*) AS count").From(blockedTable),
	}
}

func scanBlock(scanner interface{ Scan(dest ...any) error }) (DistributorBlock, error) {
	var s DistributorBlock
	var nt sql.NullTime
	if err := scanner.Scan(
		&s.ID,
		&s.DistributorID,
		&s.InitiatorID,
		&s.Reason,
		&s.Status,
		&s.BlockedAt,
		&nt, // сканим сюда
	); err != nil {
		return s, err
	}
	if nt.Valid {
		t := nt.Time
		s.CanceledAt = &t
	} else {
		s.CanceledAt = nil
	}
	return s, nil
}

func (q BlockagesQ) New() BlockagesQ {
	return NewBlockagesQ(q.db)
}

func (q BlockagesQ) Insert(ctx context.Context, input DistributorBlock) error {
	values := map[string]interface{}{
		"id":             input.ID,
		"distributor_id": input.DistributorID,
		"initiator_id":   input.InitiatorID,
		"reason":         input.Reason,
		"status":         input.Status,
		"blocked_at":     input.BlockedAt,
		//"canceled_at":    input.CanceledAt, defaults to NULL if not set
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building inserter query for table: %s input: %w", blockedTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q BlockagesQ) Get(ctx context.Context) (DistributorBlock, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return DistributorBlock{}, fmt.Errorf("building selector query for table %s: %w", blockedTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	s, err := scanBlock(row)
	if err != nil {
		return DistributorBlock{}, fmt.Errorf("scanning row for table %s: %w", blockedTable, err)
	}
	return s, nil
}

func (q BlockagesQ) Select(ctx context.Context) ([]DistributorBlock, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table %s: %w", blockedTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("executing query for table %s: %w", blockedTable, err)
	}
	defer rows.Close()

	var res []DistributorBlock
	for rows.Next() {
		s, err := scanBlock(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning row for table %s: %w", blockedTable, err)
		}
		res = append(res, s)
	}
	return res, nil
}

func (q BlockagesQ) Update(ctx context.Context) error {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", distributorsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q BlockagesQ) UpdateStatus(status string) BlockagesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q BlockagesQ) UpdateCanceledAt(canceledAt time.Time) BlockagesQ {
	q.updater = q.updater.Set("canceled_at", canceledAt)
	return q
}

func (q BlockagesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", blockedTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q BlockagesQ) FilterID(id uuid.UUID) BlockagesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q BlockagesQ) FilterDistributorID(distributorID ...uuid.UUID) BlockagesQ {
	q.selector = q.selector.Where(sq.Eq{"distributor_id": distributorID})
	q.counter = q.counter.Where(sq.Eq{"distributor_id": distributorID})
	q.updater = q.updater.Where(sq.Eq{"distributor_id": distributorID})
	q.deleter = q.deleter.Where(sq.Eq{"distributor_id": distributorID})
	return q
}

func (q BlockagesQ) FilterInitiatorID(initiatorID ...uuid.UUID) BlockagesQ {
	q.selector = q.selector.Where(sq.Eq{"initiator_id": initiatorID})
	q.counter = q.counter.Where(sq.Eq{"initiator_id": initiatorID})
	q.updater = q.updater.Where(sq.Eq{"initiator_id": initiatorID})
	q.deleter = q.deleter.Where(sq.Eq{"initiator_id": initiatorID})
	return q
}

func (q BlockagesQ) FilterStatus(status ...string) BlockagesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q BlockagesQ) OrderByBlockedAt(ascending bool) BlockagesQ {
	if ascending {
		q.selector = q.selector.OrderBy("blocked_at ASC")
	} else {
		q.selector = q.selector.OrderBy("blocked_at DESC")
	}
	return q
}

func (q BlockagesQ) Page(limit, offset uint64) BlockagesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(1) // For counting, we don't need to limit the results

	return q
}

func (q BlockagesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", blockedTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("executing count query for table %s: %w", blockedTable, err)
	}

	return count, nil
}
