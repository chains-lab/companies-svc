package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const blockedTable = "distributor_blockages"

type Blockages struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID uuid.UUID  `db:"distributor_id"`
	InitiatorID   uuid.UUID  `db:"initiator_id"`
	Reason        string     `db:"reason"`
	Status        string     `db:"status"`
	BlockedAt     time.Time  `db:"blocked_at"`
	CanceledAt    *time.Time `db:"canceled_at"`
}

type BlockQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewBlockagesQ(db *sql.DB) BlockQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return BlockQ{
		db:       db,
		selector: builder.Select("*").From(blockedTable),
		updater:  builder.Update(blockedTable),
		inserter: builder.Insert(blockedTable),
		deleter:  builder.Delete(blockedTable),
		counter:  builder.Select("COUNT(*) AS count").From(blockedTable),
	}
}

func (q BlockQ) applyConditions(conditions ...sq.Sqlizer) BlockQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)

	return q
}

func scanBlock(scanner interface{ Scan(dest ...any) error }) (Blockages, error) {
	var s Blockages
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

func (q BlockQ) New() BlockQ {
	return NewBlockagesQ(q.db)
}

func (q BlockQ) Insert(ctx context.Context, input Blockages) error {
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

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q BlockQ) Get(ctx context.Context) (Blockages, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Blockages{}, fmt.Errorf("building selector query for table %s: %w", blockedTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	s, err := scanBlock(row)
	if err != nil {
		return Blockages{}, fmt.Errorf("scanning row for table %s: %w", blockedTable, err)
	}
	return s, nil
}

func (q BlockQ) Select(ctx context.Context) ([]Blockages, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table %s: %w", blockedTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("executing query for table %s: %w", blockedTable, err)
	}
	defer rows.Close()

	var res []Blockages
	for rows.Next() {
		s, err := scanBlock(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning row for table %s: %w", blockedTable, err)
		}
		res = append(res, s)
	}
	return res, nil
}

func (q BlockQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if active, ok := input["active"]; ok {
		values["active"] = active
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", blockedTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err

}

func (q BlockQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", blockedTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q BlockQ) FilterID(id uuid.UUID) BlockQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q BlockQ) FilterDistributorID(distributorID uuid.UUID) BlockQ {
	return q.applyConditions(sq.Eq{"distributor_id": distributorID})
}

func (q BlockQ) FilterInitiatorID(initiatorID uuid.UUID) BlockQ {
	return q.applyConditions(sq.Eq{"initiator_id": initiatorID})
}

func (q BlockQ) FilterStatus(status string) BlockQ {
	return q.applyConditions(sq.Eq{"status": status})
}

func (q BlockQ) OrderByBlockedAt(ascending bool) BlockQ {
	if ascending {
		q.selector = q.selector.OrderBy("blocked_at ASC")
	} else {
		q.selector = q.selector.OrderBy("blocked_at DESC")
	}
	return q
}

func (q BlockQ) Page(limit, offset uint64) BlockQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(1) // For counting, we don't need to limit the results

	return q
}

func (q BlockQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", blockedTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("executing count query for table %s: %w", blockedTable, err)
	}

	return count, nil
}
