package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/google/uuid"
)

const employeesTable = "employees"

type Employee struct {
	UserID    uuid.UUID `db:"user_id"`
	CompanyID uuid.UUID `db:"company_id"`
	Role      string    `db:"role"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type EmployeesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewEmployeesQ(db *sql.DB) EmployeesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return EmployeesQ{
		db:       db,
		selector: builder.Select("*").From(employeesTable),
		updater:  builder.Update(employeesTable),
		inserter: builder.Insert(employeesTable),
		deleter:  builder.Delete(employeesTable),
		counter:  builder.Select("COUNT(*) AS count").From(employeesTable),
	}
}

func (q EmployeesQ) New() EmployeesQ {
	return NewEmployeesQ(q.db)
}

func (q EmployeesQ) Insert(ctx context.Context, in Employee) error {
	qry, args, err := q.inserter.
		Columns("user_id", "company_id", "role", "updated_at", "created_at").
		Values(in.UserID, in.CompanyID, in.Role, in.UpdatedAt, in.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", employeesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q EmployeesQ) Get(ctx context.Context) (Employee, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Employee{}, err
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var emp Employee
	err = row.Scan(
		&emp.UserID,
		&emp.CompanyID,
		&emp.Role,
		&emp.UpdatedAt,
		&emp.CreatedAt,
	)
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (q EmployeesQ) Select(ctx context.Context) ([]Employee, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
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

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(
			&emp.UserID,
			&emp.CompanyID,
			&emp.Role,
			&emp.UpdatedAt,
			&emp.CreatedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

func (q EmployeesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", employeesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q EmployeesQ) UpdateRole(role string) EmployeesQ {
	q.updater = q.updater.Set("role", role)
	return q
}

func (q EmployeesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for table %s: %w", employeesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q EmployeesQ) FilterUserID(userID uuid.UUID) EmployeesQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q EmployeesQ) FilterCompanyID(companyID ...uuid.UUID) EmployeesQ {
	q.selector = q.selector.Where(sq.Eq{"company_id": companyID})
	q.counter = q.counter.Where(sq.Eq{"company_id": companyID})
	q.updater = q.updater.Where(sq.Eq{"company_id": companyID})
	q.deleter = q.deleter.Where(sq.Eq{"company_id": companyID})
	return q
}

func (q EmployeesQ) FilterRole(role ...string) EmployeesQ {
	q.selector = q.selector.Where(sq.Eq{"role": role})
	q.counter = q.counter.Where(sq.Eq{"role": role})
	q.updater = q.updater.Where(sq.Eq{"role": role})
	q.deleter = q.deleter.Where(sq.Eq{"role": role})
	return q
}

func (q EmployeesQ) OrderByRole(ascend bool) EmployeesQ {
	dir := "DESC"
	if ascend {
		dir = "ASC"
	}

	caseExpr := "CASE role"
	for r, w := range enum.AllEmployeeRoles {
		caseExpr += fmt.Sprintf(" WHEN '%s' THEN %d", r, w)
	}
	caseExpr += " ELSE 0 END " + dir

	q.selector = q.selector.OrderBy(caseExpr)
	return q
}

func (q EmployeesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", employeesTable, err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q EmployeesQ) Page(limit, offset uint64) EmployeesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
