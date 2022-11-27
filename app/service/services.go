package service

import (
	"context"
	
	"github.com/google/uuid"

	"gitlab.com/m0ta/lts/app/model"
)

// UserService is a service for users
type UserService interface {
	Validate(context.Context, *model.User) error
	Get(context.Context, uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
	Update(context.Context, *model.User) (*model.User, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context) (model.Users, error)

	GetByEmail(context.Context, string) (*model.User, error)
	Exists(context.Context, string) error

	SignIn(context.Context, *model.User) (*model.User, error)
	SignUp(context.Context, *model.User) (*model.User, error)
	ChangePassword(context.Context, *model.User, string) (*model.User, error)
}

// TokenService is a service for tokens
type TokenService interface {
	Validate(context.Context, *model.Token) error
	Get(context.Context, uint64) (*model.Token, error)
	Create(context.Context, *model.Token) (*model.Token, error)
	Update(context.Context, *model.Token) (*model.Token, error)
	Delete(context.Context, uint64) error
	List(context.Context) (model.Tokens, error)

	GetByDomain(context.Context, string) (*model.Token, error)
	Exists(context.Context, string) error
	Verify(context.Context, *model.Token) (*model.ValidData, error)
}