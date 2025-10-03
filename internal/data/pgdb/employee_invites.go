package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const invitesTable = "employee_invites"

type Invite struct {
	ID         uuid.UUID     `db:"id"`
	Status     string        `db:"status"` // 'sent' | 'accepted' |
	Role       string        `db:"role"`   // enum employee_roles
	CompanyID  uuid.UUID     `db:"company_id"`
	Token      string        `db:"token"`
	UserID     uuid.NullUUID `db:"user_id"`
	AnsweredAt sql.NullTime  `db:"answered_at"`
	ExpiresAt  time.Time     `db:"expires_at"`
	CreatedAt  time.Time     `db:"created_at"`
}

type InvitesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewInvitesQ(db *sql.DB) InvitesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return InvitesQ{
		db:       db,
		selector: b.Select("*").From(invitesTable),
		inserter: b.Insert(invitesTable),
		updater:  b.Update(invitesTable),
		deleter:  b.Delete(invitesTable),
		counter:  b.Select("COUNT(*) AS count").From(invitesTable),
	}
}

func (q InvitesQ) New() InvitesQ { return NewInvitesQ(q.db) }

func (q InvitesQ) Insert(ctx context.Context, input Invite) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"status":     input.Status,
		"role":       input.Role,
		"company_id": input.CompanyID,
		"token":      input.Token,
		"expires_at": input.ExpiresAt,
	}

	if input.UserID.Valid {
		values["user_id"] = input.UserID
	}
	if input.AnsweredAt.Valid {
		values["answered_at"] = input.AnsweredAt
	}
	if !input.CreatedAt.IsZero() {
		values["created_at"] = input.CreatedAt
	}

	sqlStr, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", invitesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, sqlStr, args...)
	} else {
		_, err = q.db.ExecContext(ctx, sqlStr, args...)
	}
	return err
}

func (q InvitesQ) Get(ctx context.Context) (Invite, error) {
	sqlStr, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Invite{}, fmt.Errorf("build select %s: %w", invitesTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, sqlStr, args...)
	} else {
		row = q.db.QueryRowContext(ctx, sqlStr, args...)
	}

	var m Invite
	var userID uuid.NullUUID
	var answeredAt sql.NullTime

	if err := row.Scan(
		&m.ID,
		&m.Status,
		&m.Role,
		&m.CompanyID,
		&m.Token,
		&userID,
		&answeredAt,
		&m.ExpiresAt,
		&m.CreatedAt,
	); err != nil {
		return Invite{}, err
	}

	if userID.Valid {
		m.UserID = userID
	}
	if answeredAt.Valid {
		m.AnsweredAt = answeredAt
	}

	return m, nil
}

func (q InvitesQ) Select(ctx context.Context) ([]Invite, error) {
	sqlStr, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", invitesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, sqlStr, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, sqlStr, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Invite
	for rows.Next() {
		var m Invite
		var userID uuid.NullUUID
		var answeredAt sql.NullTime
		if err := rows.Scan(
			&m.ID,
			&m.Status,
			&m.Role,
			&m.CompanyID,
			&m.Token,
			&userID,
			&answeredAt,
			&m.ExpiresAt,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		if userID.Valid {
			m.UserID = userID
		}
		if answeredAt.Valid {
			m.AnsweredAt = answeredAt
		}
		out = append(out, m)
	}
	return out, nil
}

type UpdateInviteParams struct {
	Status     *string
	Role       *string
	CompanyID  *uuid.UUID
	UserID     *uuid.NullUUID
	AnsweredAt *sql.NullTime
	ExpiresAt  *time.Time
}

func (q InvitesQ) Update(ctx context.Context) error {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", invitesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q InvitesQ) UpdateStatus(status string) InvitesQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q InvitesQ) UpdateRole(role string) InvitesQ {
	q.updater = q.updater.Set("role", role)
	return q
}

func (q InvitesQ) UpdateCompanyID(companyID uuid.UUID) InvitesQ {
	q.updater = q.updater.Set("company_id", companyID)
	return q
}

func (q InvitesQ) UpdateUserID(userID uuid.UUID) InvitesQ {
	q.updater = q.updater.Set("user_id", userID)
	return q
}

func (q InvitesQ) UpdateAnsweredAt(answeredAt time.Time) InvitesQ {
	q.updater = q.updater.Set("answered_at", answeredAt)
	return q
}

func (q InvitesQ) UpdateExpiresAt(expiresAt time.Time) InvitesQ {
	q.updater = q.updater.Set("expires_at", expiresAt)
	return q
}

func (q InvitesQ) Delete(ctx context.Context) error {
	sqlStr, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", invitesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, sqlStr, args...)
	} else {
		_, err = q.db.ExecContext(ctx, sqlStr, args...)
	}
	return err
}

func (q InvitesQ) FilterID(id uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q InvitesQ) FiltercompanyID(companyID uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"company_id": companyID})
	q.updater = q.updater.Where(sq.Eq{"company_id": companyID})
	q.deleter = q.deleter.Where(sq.Eq{"company_id": companyID})
	q.counter = q.counter.Where(sq.Eq{"company_id": companyID})
	return q
}

func (q InvitesQ) FilterUserID(userID uuid.UUID) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q InvitesQ) FilterStatus(status ...string) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	return q
}

func (q InvitesQ) FilterRole(role ...string) InvitesQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	return q
}

func (q InvitesQ) FilterExpiresBefore(t time.Time) InvitesQ {
	q.selector = q.selector.Where(sq.LtOrEq{"expires_at": t})
	q.updater = q.updater.Where(sq.LtOrEq{"expires_at": t})
	q.deleter = q.deleter.Where(sq.LtOrEq{"expires_at": t})
	q.counter = q.counter.Where(sq.LtOrEq{"expires_at": t})
	return q
}

func (q InvitesQ) FilterExpiresAfter(t time.Time) InvitesQ {
	q.selector = q.selector.Where(sq.Gt{"expires_at": t})
	q.updater = q.updater.Where(sq.Gt{"expires_at": t})
	q.deleter = q.deleter.Where(sq.Gt{"expires_at": t})
	q.counter = q.counter.Where(sq.Gt{"expires_at": t})
	return q
}

func (q InvitesQ) FilterAnswered(answered bool) InvitesQ {
	if answered {
		q.selector = q.selector.Where("answered_at IS NOT NULL")
		q.updater = q.updater.Where("answered_at IS NOT NULL")
		q.deleter = q.deleter.Where("answered_at IS NOT NULL")
		q.counter = q.counter.Where("answered_at IS NOT NULL")
	} else {
		q.selector = q.selector.Where("answered_at IS NULL")
		q.updater = q.updater.Where("answered_at IS NULL")
		q.deleter = q.deleter.Where("answered_at IS NULL")
		q.counter = q.counter.Where("answered_at IS NULL")
	}
	return q
}

func (q InvitesQ) FilterCreatedBetween(from, to time.Time) InvitesQ {
	q.selector = q.selector.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.updater = q.updater.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.deleter = q.deleter.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	q.counter = q.counter.Where(sq.And{
		sq.GtOrEq{"created_at": from},
		sq.LtOrEq{"created_at": to},
	})
	return q
}

func (q InvitesQ) OrderByCreatedAt(asc bool) InvitesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("created_at " + dir)
	return q
}

func (q InvitesQ) OrderByUpdatedAt(asc bool) InvitesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("updated_at " + dir)
	return q
}

func (q InvitesQ) OrderByExpiresAt(asc bool) InvitesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	q.selector = q.selector.OrderBy("expires_at " + dir)
	return q
}

func (q InvitesQ) Count(ctx context.Context) (uint64, error) {
	sqlStr, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", invitesTable, err)
	}

	var n uint64
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, sqlStr, args...)
	} else {
		row = q.db.QueryRowContext(ctx, sqlStr, args...)
	}
	if err := row.Scan(&n); err != nil {
		return 0, fmt.Errorf("scan count %s: %w", invitesTable, err)
	}
	return n, nil
}

func (q InvitesQ) Page(limit, offset uint64) InvitesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
