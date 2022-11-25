package bun

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"gitlab.com/m0ta/lts/app/config"
	"gitlab.com/m0ta/lts/app/utils"
)

// Timeout is a Postgres timeout
const Timeout = 5

// DB is a shortcut structure to a Postgres DB
type DB struct {
	*bun.DB
}

// Dial creates new database connection to postgres
func Dial() (*DB, error) {
	cfg := config.Get()
	if cfg.PgURL == "" {
		return nil, utils.ErrorNew("No URL to connect Postgre (bun/Dial)")
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.PgURL)))
	db := bun.NewDB(sqlDB, pgdialect.New())

	//for furture debug
	db.AddQueryHook(bundebug.NewQueryHook())

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}