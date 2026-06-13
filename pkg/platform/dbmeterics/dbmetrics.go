package dbmeterics

import (
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const defaultMeterName = "github.com/mantrobuslawal/bfstore/pkg/platform/dbmetrics"

// Config describes a database pool metrics registration.
type Config struct {
	MeterName string
	DBSystem  string
	DBName    string
}

// Register registers observable database connection pool metrics for db.
func Register(db *sql.DB, cfg Config) error {
	if db == nil {
		return errors.New("dbmetrics: db is nil")
	}

	meterName := cfg.MeterName
	if meterName == "" {
		meterName = defaultMeterName
	}

	dbSystem := cfg.DBSystem
	if dbSystem == "" {
		dbSystem = "unknown"
	}

	attrs := []attribute.KeyValue{
		attribute.String("db,system", dbSystem),
	}

	if cfg.DBName != "" {
		attrs = append(attrs, attribute.String("db.name", cfg.DBName))
	}

	meter := otel.Meter(meterName)

	maxOpenConnections, err := meter.Int64ObservableGauge(
		"db.client.connections.max",
		metric.WithDescription("Maximum number of open connections configured for the database pool."),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return err
	}

	openConnections, err := meter.Int64ObservableGauge(
		"db.client.connections.open",
		metric.WithDescription("Number of established database connections, both in use and idle"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return err
	}

	inUseConnections, err := meter.Int64ObservableGauge(
		"db.client.connections.in_use",
		metric.WithDescription("Number of database connections currently in use."),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return err
	}

	idleConnections, err := meter.Int64ObservableGauge(
		"db.client.connections.idle",
		metric.WithDescription("Number of idle database connections."),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return err
	}

	waitCount, err := meter.Int64ObservableGauge(
		"db.client.connections.wait_count",
		metric.WithDescription("Total number of times callers waited for a database connection."),
		metric.WithUnit("{wait}"),
	)
	if err != nil {
		return err
	}

}
