package service

import (
	"context"
	"log"
	"os"

	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

type Logger struct {
	Info 	*log.Logger
	Error 	*log.Logger
}

// Manager is just a collection of all services we have in the project
type Manager struct {
	Logger		Logger
	User        UserService
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, utils.ErrorNew("No store provided")
	}
	return &Manager{
		Logger: 	NewLoggers(),
		User:       NewUserWebService(ctx, store),
	}, nil
}

func NewLoggers() Logger {
	var l Logger
	l.Info 	= log.New(os.Stdout, "INFO\t", 	log.Ldate|log.Ltime|log.LUTC)
	l.Error = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC)
	return l
}