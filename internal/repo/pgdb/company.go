package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const companiesTable = "companies"

type Company struct {
	ID        uuid.UUID `db:"id"`
	Icon      string    `db:"icon"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type CompaniesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewcompaniesQ(db *sql.DB) CompaniesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return CompaniesQ{
		db:       db,
		selector: builder.Select("*").From(companiesTable),
		updater:  builder.Update(companiesTable),
		inserter: builder.Insert(companiesTable),
		deleter:  builder.Delete(companiesTable),
		counter:  builder.Select("COUNT(*) AS count").From(companiesTable),
	}
}

func (q CompaniesQ) New() CompaniesQ {
	return NewcompaniesQ(q.db)
}

func (q CompaniesQ) Insert(ctx context.Context, in Company) error {
	qry, args, err := q.inserter.
		Columns("id", "icon", "name", "status", "updated_at", "created_at").
		Values(in.ID, in.Icon, in.Name, in.Status, in.UpdatedAt, in.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", companiesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q CompaniesQ) Get(ctx context.Context) (Company, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Company{}, fmt.Errorf("building selector query for table: %s: %w", companiesTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var d Company
	err = row.Scan(
		&d.ID,
		&d.Icon,
		&d.Name,
		&d.Status,
		&d.UpdatedAt,
		&d.CreatedAt,
	)
	if err != nil {
		return Company{}, fmt.Errorf("scanning row for table: %s: %w", companiesTable, err)
	}

	return d, nil
}

func (q CompaniesQ) Select(ctx context.Context) ([]Company, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", companiesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var d Company
		if err := rows.Scan(
			&d.ID,
			&d.Icon,
			&d.Name,
			&d.Status,
			&d.UpdatedAt,
			&d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning row for table: %s: %w", companiesTable, err)
		}
		companies = append(companies, d)
	}

	return companies, nil
}

func (q CompaniesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", companiesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q CompaniesQ) UpdateName(name string) CompaniesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q CompaniesQ) UpdateIcon(icon string) CompaniesQ {
	q.updater = q.updater.Set("icon", icon)
	return q
}

func (q CompaniesQ) UpdateStatus(status string) CompaniesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q CompaniesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", companiesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CompaniesQ) FilterID(id uuid.UUID) CompaniesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	return q
}

func (q CompaniesQ) FilterStatus(status ...string) CompaniesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	return q
}

func (q CompaniesQ) FilterLikeName(name string) CompaniesQ {
	cond := sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", name))
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q CompaniesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building counter query for table: %s: %w", companiesTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("scanning count for table: %s: %w", companiesTable, err)
	}

	return count, nil
}

func (q CompaniesQ) Page(limit, offset uint64) CompaniesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q CompaniesQ) OrderByName(ascend bool) CompaniesQ {
	if ascend {
		q.selector = q.selector.OrderBy("name ASC")
	} else {
		q.selector = q.selector.OrderBy("name DESC")
	}
	return q
}

func (q CompaniesQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
