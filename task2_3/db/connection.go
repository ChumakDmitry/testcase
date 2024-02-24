package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"main/config"
)

type Postgres struct {
	db *pgxpool.Pool
}

var pgInstance *Postgres

func InitPG(ctx context.Context, cfg config.Config) (*Postgres, error) {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := pgxpool.New(ctx, connectStr)
	if err != nil {
		log.Fatalf("error to create connection pool %s", err)
	}

	pgInstance = &Postgres{db}

	if err := pgInstance.db.Ping(ctx); err != nil {
		panic(err)
	}

	return pgInstance, nil
}
