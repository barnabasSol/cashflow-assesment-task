package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	ConnPool *pgxpool.Pool
}

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
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	const DATABASE_URL = "postgres://postgres:strongpassword@localhost:5432/app_db?sslmode=disable"

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}
