package store

import (
	"context"
	"time"
	"log"

	"github.com/pkg/errors"

	"gitlab.com/m0ta/lts/app/store/bun"
)

// Store ...
type Store struct {
	Bun    		*bun.DB       // for KeepAlivePg (see below)
	User 		UserRepo
}

// New creates new store
func New(ctx context.Context) (*Store, error) {
	//cfg := config.Get()

	// connect to Postgres
	bunDB, err := bun.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "bun.Dial failed")
	}
	
	var store Store

	// Init Postgres repositories
	if bunDB != nil {
		store.Bun = bunDB
		go store.KeepAliveBun()
		store.User 		= bun.NewUserRepo(bunDB)
	}

	return &store, nil
}

// KeepAlivePollPeriod is a Pg/MySQL keepalive check time period
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
		log.Println("[store.KeepAliveBun] Lost PostgreSQL connection. Restoring...")
		store.Bun, err = bun.Dial()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		log.Println("[store.KeepAliveBun] PostgreSQL reconnected")
	}
}