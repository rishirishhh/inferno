package repository

import (
	db "auth-service/internal/db/generated"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

func New(databaseURL string) (*Repository, error) {
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Connection Pool Configuration
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxIdleTime(15 * time.Minute)
	dbConn.SetConnMaxLifetime(1 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.PingContext(ctx); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("ping database")
	}

	return &Repository{
		db:      dbConn,
		queries: db.New(dbConn),
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) WithTx(
	ctx context.Context,
	fn func(*db.Queries) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := r.queries.WithTx(tx)

	if err := fn(q); err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf(
				"transaction failed : %v, rollback failed: %v",
				err,
				rollBackErr,
			)
		}
		return err
	}

	return tx.Commit()
}

func (r *Repository) Queries() *db.Queries {
	return r.queries
}
