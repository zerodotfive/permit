package types

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Config   *Config
	Logger   *log.Logger
	Postgres *pgxpool.Pool
}
