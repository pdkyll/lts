package store

import (
	"context"

	"github.com/google/uuid"

	"gitlab.com/m0ta/lts/app/model"
)

// UserRepo is a store for users
type UserRepo interface {
	GetByEmail(context.Context, string) (*model.User, error)

	Get(context.Context, uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
	Update(context.Context, *model.User) (*model.User, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context) (model.Users, error)
}