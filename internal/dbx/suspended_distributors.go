package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const suspendedDistributorsTable = "suspended_distributors"

type SuspendedDistributor struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID uuid.UUID  `db:"distributor_id"`
	InitiatorID   uuid.UUID  `db:"initiator_id"`
	Reason        string     `db:"reason"`
	Active        bool       `db:"active"`
	SuspendedAt   time.Time  `db:"suspended_at"`
	CanceledAt    *time.Time `db:"canceled_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type SuspendedQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewSuspendedDistributorsQ(db *sql.DB) SuspendedQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return SuspendedQ{
		db:       db,
		selector: builder.Select("*").From(suspendedDistributorsTable),
		updater:  builder.Update(suspendedDistributorsTable),
		inserter: builder.Insert(suspendedDistributorsTable),
		deleter:  builder.Delete(suspendedDistributorsTable),
		counter:  builder.Select("COUNT(*) AS count").From(suspendedDistributorsTable),
	}
}

func (q SuspendedQ) applyConditions(conditions ...sq.Sqlizer) SuspendedQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)

	return q
}

func scanSuspended(scanner interface{ Scan(dest ...any) error }) (SuspendedDistributor, error) {
	var s SuspendedDistributor
	var nt sql.NullTime
	if err := scanner.Scan(
		&s.ID,
		&s.DistributorID,
		&s.InitiatorID,
		&s.Reason,
		&s.Active,
		&s.SuspendedAt,
		&nt, // сканим сюда
		&s.CreatedAt,
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

func (q SuspendedQ) New() SuspendedQ {
	return NewSuspendedDistributorsQ(q.db)
}

func (q SuspendedQ) Insert(ctx context.Context, input SuspendedDistributor) error {
	values := map[string]interface{}{
		"id":             input.ID,
		"distributor_id": input.DistributorID,
		"initiator_id":   input.InitiatorID,
		"reason":         input.Reason,
		"active":         input.Active,
		"suspended_at":   input.SuspendedAt,
		//"canceled_at":    input.CanceledAt, defaults to NULL if not set
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building inserter query for table: %s input: %w", suspendedDistributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q SuspendedQ) Get(ctx context.Context) (SuspendedDistributor, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return SuspendedDistributor{}, fmt.Errorf("building selector query for table %s: %w", suspendedDistributorsTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	s, err := scanSuspended(row)
	if err != nil {
		return SuspendedDistributor{}, fmt.Errorf("scanning row for table %s: %w", suspendedDistributorsTable, err)
	}
	return s, nil
}

func (q SuspendedQ) Select(ctx context.Context) ([]SuspendedDistributor, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table %s: %w", suspendedDistributorsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("executing query for table %s: %w", suspendedDistributorsTable, err)
	}
	defer rows.Close()

	var res []SuspendedDistributor
	for rows.Next() {
		s, err := scanSuspended(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning row for table %s: %w", suspendedDistributorsTable, err)
		}
		res = append(res, s)
	}
	return res, nil
}

func (q SuspendedQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if active, ok := input["active"]; ok {
		values["active"] = active
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", suspendedDistributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err

}

func (q SuspendedQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", suspendedDistributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q SuspendedQ) FilterID(id uuid.UUID) SuspendedQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q SuspendedQ) FilterDistributorID(distributorID uuid.UUID) SuspendedQ {
	return q.applyConditions(sq.Eq{"distributor_id": distributorID})
}

func (q SuspendedQ) FilterInitiatorID(initiatorID uuid.UUID) SuspendedQ {
	return q.applyConditions(sq.Eq{"initiator_id": initiatorID})
}

func (q SuspendedQ) FilterActive(active bool) SuspendedQ {
	return q.applyConditions(sq.Eq{"active": active})
}

func (q SuspendedQ) OrderBySuspendedAt(ascending bool) SuspendedQ {
	if ascending {
		q.selector = q.selector.OrderBy("suspended_at ASC")
	} else {
		q.selector = q.selector.OrderBy("suspended_at DESC")
	}
	return q
}

func (q SuspendedQ) Page(limit, offset uint64) SuspendedQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(1) // For counting, we don't need to limit the results

	return q
}

func (q SuspendedQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", suspendedDistributorsTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("executing count query for table %s: %w", suspendedDistributorsTable, err)
	}

	return count, nil
}
