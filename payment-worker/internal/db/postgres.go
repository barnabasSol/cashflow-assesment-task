package db

import (
	"context"
	"log"
	"payment-worker/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	ConnPool *pgxpool.Pool
}

const (
	defaultMaxConns          = int32(10)
	defaultMinConns          = int32(2)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = 30 * time.Minute
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = 5 * time.Second
)

func InitPostgres(ctx context.Context) *Postgres {
	connPool, err := pgxpool.NewWithConfig(ctx, Config())
	if err != nil {
		log.Fatal("Error while creating connection to the database")
	}
	err = connPool.Ping(ctx)
	if err != nil {
		log.Fatal("Error while pinging the database")
	}
	return &Postgres{
		ConnPool: connPool,
	}
}

func Config() *pgxpool.Config {
	db_url := config.GetSecret("DB_URL")

	cfg, err := pgxpool.ParseConfig(db_url)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	cfg.MaxConns = config.GetEnvInt32(
		"DB_MAX_CONNS",
		defaultMaxConns,
	)

	cfg.MinConns = config.GetEnvInt32(
		"DB_MIN_CONNS",
		defaultMinConns,
	)

	cfg.MaxConnLifetime = config.GetEnvDuration(
		"DB_MAX_CONN_LIFETIME",
		defaultMaxConnLifetime,
	)

	cfg.MaxConnIdleTime = config.GetEnvDuration(
		"DB_MAX_CONN_IDLE_TIME",
		defaultMaxConnIdleTime,
	)

	cfg.HealthCheckPeriod = config.GetEnvDuration(
		"DB_HEALTH_CHECK_PERIOD",
		defaultHealthCheckPeriod,
	)

	cfg.ConnConfig.ConnectTimeout = config.GetEnvDuration(
		"DB_CONNECT_TIMEOUT",
		defaultConnectTimeout,
	)

	return cfg
}
