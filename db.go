package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	DBPool *pgxpool.Pool
}

// NewDatabase initializes the read and write database pools.
func NewDatabase(ctx context.Context, DSN string, logger *slog.Logger) (*Database, error) {
	logger.Info("setting up database", "DSN", DSN)
	DBPool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		logger.Error("failed to setup read database", "error", err)
		return nil, err
	}

	err = DBPool.Ping(ctx)
	if err != nil {
		logger.Error("failed to ping database", "error", err)
		return nil, err
	}
	logger.Info("database setup complete", "DSN", DSN)
	return &Database{
		DBPool: DBPool,
	}, nil
}

// WithTransaction manages database transactions, ensuring proper commit or rollback.
// If the provided function (fn) returns an error, the transaction is rolled back; otherwise, it is committed.
func WithTransaction[T any](ctx context.Context, dbConn *Database, fn func(pgx.Tx) (T, error)) (T, error) {
	tx, err := dbConn.DBPool.Begin(ctx)
	if err != nil {
		return *new(T), fmt.Errorf("beginning transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	result, err := fn(tx)
	if err != nil {
		return *new(T), err
	}

	if err = tx.Commit(ctx); err != nil {
		return *new(T), fmt.Errorf("committing transaction: %w", err)
	}
	return result, nil
}

// Close closes the database pools.
func (db *Database) Close() {
	db.DBPool.Close()
}
