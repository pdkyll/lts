package service

import (
	"context"
	
	"go.uber.org/zap"

	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	Logger	*zap.Logger
	User    UserService
	Token	TokenService
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *store.Store, logger *zap.Logger) (*Manager, error) {
	if store == nil {
		return nil, utils.ErrorNew("No store provided")
	}
	return &Manager{
		Logger:	logger,
		User:   NewUserWebService(ctx, store, logger),
		Token:  NewTokenWebService(ctx, store, logger),
	}, nil
}