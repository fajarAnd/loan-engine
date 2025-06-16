package database

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pgx-contrib/pgxotel"
)

var (
	db *pgxpool.Pool
)

func InitDB(connStr string) error {
	pgxPoolConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing db config: %w", err)
	}

	pgxPoolConf.ConnConfig.Tracer = &pgxotel.QueryTracer{
		Name: "loan-service",
	}

	db, err = pgxpool.NewWithConfig(context.Background(), pgxPoolConf)
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}

	return nil
}

func GetConn() *pgxpool.Pool {
	return db
}

var ProviderSet = wire.NewSet(
	GetConn,
	wire.Bind(new(Querier), new(*pgxpool.Pool)),
	wire.Bind(new(Executor), new(*pgxpool.Pool)),
	wire.Bind(new(Tx), new(*pgxpool.Pool)),
)
