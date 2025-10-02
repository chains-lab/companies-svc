package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const distributorsTable = "distributors"

type Distributor struct {
	ID        uuid.UUID `db:"id"`
	Icon      string    `db:"icon"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type DistributorsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewDistributorsQ(db *sql.DB) DistributorsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return DistributorsQ{
		db:       db,
		selector: builder.Select("*").From(distributorsTable),
		updater:  builder.Update(distributorsTable),
		inserter: builder.Insert(distributorsTable),
		deleter:  builder.Delete(distributorsTable),
		counter:  builder.Select("COUNT(*) AS count").From(distributorsTable),
	}
}

func (q DistributorsQ) New() DistributorsQ {
	return NewDistributorsQ(q.db)
}

func (q DistributorsQ) Insert(ctx context.Context, in Distributor) error {
	qry, args, err := q.inserter.
		Columns("id", "icon", "name", "status", "updated_at", "created_at").
		Values(in.ID, in.Icon, in.Name, in.Status, in.UpdatedAt, in.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", distributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q DistributorsQ) Get(ctx context.Context) (Distributor, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Distributor{}, fmt.Errorf("building selector query for table: %s: %w", distributorsTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var d Distributor
	err = row.Scan(
		&d.ID,
		&d.Icon,
		&d.Name,
		&d.Status,
		&d.UpdatedAt,
		&d.CreatedAt,
	)
	if err != nil {
		return Distributor{}, fmt.Errorf("scanning row for table: %s: %w", distributorsTable, err)
	}

	return d, nil
}

func (q DistributorsQ) Select(ctx context.Context) ([]Distributor, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", distributorsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributors []Distributor
	for rows.Next() {
		var d Distributor
		if err := rows.Scan(
			&d.ID,
			&d.Icon,
			&d.Name,
			&d.Status,
			&d.UpdatedAt,
			&d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning row for table: %s: %w", distributorsTable, err)
		}
		distributors = append(distributors, d)
	}

	return distributors, nil
}

func (q DistributorsQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if icon, ok := input["icon"]; ok {
		values["icon"] = icon
	}
	if name, ok := input["name"]; ok {
		values["name"] = name
	}
	if status, ok := input["status"]; ok {
		values["status"] = status
	}
	if updatedAt, ok := input["updated_at"]; ok {
		values["updated_at"] = updatedAt
	} else {
		values["updated_at"] = time.Now().UTC()
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", distributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q DistributorsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", distributorsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q DistributorsQ) FilterID(id uuid.UUID) DistributorsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q DistributorsQ) FilterStatus(status ...string) DistributorsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q DistributorsQ) LikeName(name string) DistributorsQ {
	cond := sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", name))
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q DistributorsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building counter query for table: %s: %w", distributorsTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("scanning count for table: %s: %w", distributorsTable, err)
	}

	return count, nil
}

func (q DistributorsQ) Page(limit, offset uint64) DistributorsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q DistributorsQ) OrderByName(ascend bool) DistributorsQ {
	if ascend {
		q.selector = q.selector.OrderBy("name ASC")
	} else {
		q.selector = q.selector.OrderBy("name DESC")
	}
	return q
}

func (q DistributorsQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
