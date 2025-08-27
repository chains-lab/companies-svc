package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/google/uuid"
)

const employeesTable = "employees"

type Employee struct {
	UserID        uuid.UUID `db:"user_id"`
	DistributorID uuid.UUID `db:"distributor_id"`
	Role          string    `db:"role"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
}

type EmployeeQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewEmployeesQ(db *sql.DB) EmployeeQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return EmployeeQ{
		db:       db,
		selector: builder.Select("*").From(employeesTable),
		updater:  builder.Update(employeesTable),
		inserter: builder.Insert(employeesTable),
		deleter:  builder.Delete(employeesTable),
		counter:  builder.Select("COUNT(*) AS count").From(employeesTable),
	}
}

func (q EmployeeQ) applyConditions(conditions ...sq.Sqlizer) EmployeeQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)

	return q
}

func (q EmployeeQ) New() EmployeeQ {
	return NewEmployeesQ(q.db)
}

func (q EmployeeQ) Insert(ctx context.Context, input Employee) error {
	values := map[string]interface{}{
		"user_id":        input.UserID,
		"distributor_id": input.DistributorID,
		"role":           input.Role,
		"updated_at":     input.UpdatedAt,
		"created_at":     input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return err
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q EmployeeQ) Get(ctx context.Context) (Employee, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Employee{}, err
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var emp Employee
	err = row.Scan(
		&emp.UserID,
		&emp.DistributorID,
		&emp.Role,
		&emp.UpdatedAt,
		&emp.CreatedAt,
	)
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (q EmployeeQ) Select(ctx context.Context) ([]Employee, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
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
			&emp.DistributorID,
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

func (q EmployeeQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if role, ok := input["role"]; ok {
		values["role"] = role
	}
	if updatedAt, ok := input["updated_at"]; ok {
		values["updated_at"] = updatedAt
	} else {
		values["updated_at"] = time.Now().UTC()
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for table %s: %w", employeesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q EmployeeQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for table %s: %w", employeesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q EmployeeQ) FilterUserID(userID uuid.UUID) EmployeeQ {
	return q.applyConditions(sq.Eq{"user_id": userID})
}

func (q EmployeeQ) FilterDistributorID(distributorID ...uuid.UUID) EmployeeQ {
	return q.applyConditions(sq.Eq{"distributor_id": distributorID})
}

func (q EmployeeQ) FilterRole(role ...string) EmployeeQ {
	return q.applyConditions(sq.Eq{"role": role})
}

func (q EmployeeQ) OrderByRole(ascend bool) EmployeeQ {
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

func (q EmployeeQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", employeesTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q EmployeeQ) Page(limit, offset uint64) EmployeeQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q EmployeeQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

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
