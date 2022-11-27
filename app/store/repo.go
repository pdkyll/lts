package store

import (
	"context"

	"github.com/google/uuid"

	"gitlab.com/m0ta/lts/app/model"
)

// UserRepo is a store for users
type UserRepo interface {
	Get(context.Context, uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
	Update(context.Context, *model.User) (*model.User, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context) (model.Users, error)

	GetByEmail(context.Context, string) (*model.User, error)
	Exists(context.Context, string) (bool, error)
}

// TokenRepo is a store for tokens
type TokenRepo interface {
	Get(context.Context, uint64) (*model.Token, error)
	Create(context.Context, *model.Token) (*model.Token, error)
	Update(context.Context, *model.Token) (*model.Token, error)
	Delete(context.Context, uint64) error
	List(context.Context) (model.Tokens, error)

	GetByDomain(context.Context, string) (*model.Token, error)
	Exists(context.Context, string) (bool, error)
}