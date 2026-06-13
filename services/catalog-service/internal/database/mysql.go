package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/mantrobuslawal/bfstore/services/catalog-service/internal/config"
)

const instrumentedMySQLDriverName = "bfstore-mysql-otel"

// Open creates a MySQL database handle.
//
// The returned *sql.DB is a connection pool, not a single connection.
func Open(cfg config.DatabaseConfig) (*sql.DB, error) {
	driverName, err := otelsql.Register(instrumentedMySQLDriverName)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName, cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

// Ping verifies database connectivity.
func Ping(ctx context.Context, db *sql.DB) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.PingContext(pingCtx)
}
