package store

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/m0ta/lts/app/store/bun"
)

// Store ...
type Store struct {
	Bun    		*bun.DB
	Logger 		*zap.Logger
	User 		UserRepo
	Token 		TokenRepo
}

// New creates new store
func New(ctx context.Context, logger *zap.Logger) (*Store, error) {
	// Connect to Postgres
	bunDB, err := bun.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "bun.Dial failed")
	}
	
	var store Store

	// Init Postgres repositories
	if bunDB != nil {
		store.Bun 	= bunDB
		store.Logger= logger
		go store.KeepAliveBun()
		store.User 	= bun.NewUserRepo(bunDB, logger)
		store.Token	= bun.NewTokenRepo(bunDB, logger)
	}

	return &store, nil
}

// KeepAlivePollPeriod is a DB keepalive check time period
const KeepAlivePollPeriod = 30

// KeepAliveBun makes sure PostgreSQL is alive and reconnects if needed
func (store *Store) KeepAliveBun() {
	var err error
	for {
		// Check if PostgreSQL is alive every KeepAlivePollPeriod seconds
		time.Sleep(time.Second * KeepAlivePollPeriod)
		lostConnect := false
		if store.Bun == nil {
			lostConnect = true
		} else if err = store.Bun.Ping(); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		store.Logger.Warn("[store.KeepAliveBun] Lost PostgreSQL connection. Restoring...")
		store.Bun, err = bun.Dial()
		if err != nil {
			store.Logger.Fatal(err.Error())
			continue
		}
		store.Logger.Warn("[store.KeepAliveBun] PostgreSQL reconnected")
	}
}