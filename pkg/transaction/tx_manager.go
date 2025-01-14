package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Manager defines the transaction manager interface.
//
//go:generate mockery --name=Manager --outpkg transactionmock --output transactionmock --case underscore
type Manager interface {
	WithTx(c context.Context, callback func(tx pgx.Tx) error) error
}

type manager struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Manager {
	return &manager{db: db}
}

func (m *manager) WithTx(c context.Context, callback func(tx pgx.Tx) error) error {
	// Start a transaction
	tx, err := m.db.Begin(c)
	if err != nil {
		return err
	}

	// Defer a function to handle the transaction based on success or failure
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(c)
			panic(p) // Re-panic after rolling back
		} else if err != nil {
			tx.Rollback(c)
		} else {
			// Commit the transaction if everything is successful
			if err := tx.Commit(c); err != nil {
				tx.Rollback(c)
			}
		}
	}()

	// Execute the callback function with the transaction
	if err := callback(tx); err != nil {
		return err
	}

	return nil
}
