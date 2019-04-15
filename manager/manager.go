package manager

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type Queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func Transaction(ctx context.Context, db *sql.DB, f func(ctx context.Context, tx Queryer) error) error {
	const maxAttempts = 2
	var err error
	for i := 0; i < maxAttempts; i++ {
		err = transaction(ctx, db, f)
		if err == nil {
			return nil
		}
		pqerr, ok := err.(*pq.Error)
		if !ok || pqerr.Message != "could not serialize access due to read/write dependencies among transactions" {
			return err
		}
	}
	return err
}

func transaction(ctx context.Context, db *sql.DB, f func(ctx context.Context, tx Queryer) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil {
			log.Print(err)
		}
	}()

	if err := f(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
