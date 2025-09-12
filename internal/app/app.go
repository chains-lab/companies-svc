package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/entities/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/entities/employee"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/distributors-svc/internal/dbx"
)

type App struct {
	distributor distributor.Distributor
	employee    employee.Employee

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		distributor: distributor.NewDistributor(pg),
		employee:    employee.NewEmployee(pg, cfg),

		db: pg,
	}, nil
}

func (a App) transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, dbx.TxKey, tx)

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
