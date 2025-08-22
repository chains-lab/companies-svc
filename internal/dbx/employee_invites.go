package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const employeeInvitesTable = "employee_invites"

type Invite struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID uuid.UUID  `db:"distributor_id"`
	UserID        uuid.UUID  `db:"user_id"`
	InvitedBy     uuid.UUID  `db:"invited_by"`
	Role          string     `db:"role"` // enum employee_roles
	Status        string     `db:"status"`
	AnsweredAt    *time.Time `db:"answered_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type InviteQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewInvitesQ(db *sql.DB) InviteQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return InviteQ{
		db:       db,
		selector: builder.Select("*").From(employeeInvitesTable),
		updater:  builder.Update(employeeInvitesTable),
		inserter: builder.Insert(employeeInvitesTable),
		deleter:  builder.Delete(employeeInvitesTable),
		counter:  builder.Select("COUNT(*) AS count").From(employeeInvitesTable),
	}
}

func (q InviteQ) New() InviteQ {
	return NewInvitesQ(q.db)
}

func (q InviteQ) applyConditions(conds ...sq.Sqlizer) InviteQ {
	q.selector = q.selector.Where(conds)
	q.counter = q.counter.Where(conds)
	q.updater = q.updater.Where(conds)
	q.deleter = q.deleter.Where(conds)
	return q
}

func scanInvitation(scanner interface{ Scan(dest ...any) error }) (Invite, error) {
	var inv Invite
	var nt sql.NullTime

	err := scanner.Scan(
		&inv.ID,
		&inv.DistributorID,
		&inv.UserID,
		&inv.InvitedBy,
		&inv.Role,
		&inv.Status,
		&nt,
		&inv.CreatedAt,
	)
	if err != nil {
		return inv, err
	}
	if nt.Valid {
		t := nt.Time
		inv.AnsweredAt = &t
	}
	return inv, nil
}

func (q InviteQ) Insert(ctx context.Context, input Invite) error {
	values := map[string]interface{}{
		"id":             input.ID,
		"distributor_id": input.DistributorID,
		"user_id":        input.UserID,
		"invited_by":     input.InvitedBy,
		"role":           input.Role,
		"status":         input.Status,
		"answered_at":    input.AnsweredAt,
		"created_at":     input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for table %s: %w", employeeInvitesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q InviteQ) Get(ctx context.Context) (Invite, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("building select query for %s: %w", employeeInvitesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	inv, err := scanInvitation(row)
	if err != nil {
		return Invite{}, fmt.Errorf("scanning row for %s: %w", employeeInvitesTable, err)
	}
	return inv, nil
}

func (q InviteQ) Select(ctx context.Context) ([]Invite, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", employeeInvitesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("executing select for %s: %w", employeeInvitesTable, err)
	}
	defer rows.Close()

	var res []Invite
	for rows.Next() {
		inv, err := scanInvitation(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning row for %s: %w", employeeInvitesTable, err)
		}
		res = append(res, inv)
	}
	return res, nil
}

func (q InviteQ) Update(ctx context.Context, input map[string]any) error {
	query, args, err := q.updater.SetMap(input).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", employeeInvitesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q InviteQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", employeeInvitesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q InviteQ) FilterID(id uuid.UUID) InviteQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q InviteQ) FilterDistributorID(distributorID uuid.UUID) InviteQ {
	return q.applyConditions(sq.Eq{"distributor_id": distributorID})
}

func (q InviteQ) FilterUserID(userID uuid.UUID) InviteQ {
	return q.applyConditions(sq.Eq{"user_id": userID})
}

func (q InviteQ) FilterInvitedBy(userID uuid.UUID) InviteQ {
	return q.applyConditions(sq.Eq{"invited_by": userID})
}

func (q InviteQ) FilterRole(role string) InviteQ {
	return q.applyConditions(sq.Eq{"role": role})
}

func (q InviteQ) FilterStatus(status string) InviteQ {
	return q.applyConditions(sq.Eq{"status": status})
}

func (q InviteQ) OrderByCreatedAt(asc bool) InviteQ {
	if asc {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}
	return q
}

func (q InviteQ) Page(limit, offset uint64) InviteQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(1) // count не нужно ограничивать выборкой
	return q
}

func (q InviteQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", employeeInvitesTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, fmt.Errorf("executing count for %s: %w", employeeInvitesTable, err)
	}
	return count, nil
}
