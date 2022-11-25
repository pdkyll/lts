package service

import (
	"context"
	
	"github.com/google/uuid"
	
	"gitlab.com/m0ta/lts/app/model"
	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

// UserWebService ...
type UserWebService struct {
	ctx   context.Context
	store *store.Store
}

// NewUserWebService creates a new user web service
func NewUserWebService(ctx context.Context, store *store.Store) *UserWebService {
	return &UserWebService{
		ctx:   ctx,
		store: store,
	}
}

// Get ...
func (svc *UserWebService) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := svc.store.User.Get(ctx, userID)
	if err != nil {
		return nil, utils.ErrorWrap(err, "svc.user.Get")
	}
	
	return user, nil
}

// Create ...
func (svc UserWebService) Create(ctx context.Context, user *model.User) (*model.User, error) {

	createdUser, err := svc.store.User.Create(ctx, user)
	if err != nil {
		return nil, utils.ErrorWrap(err, "svc.user.Create error")
	}

	return createdUser, nil
}

// Update ...
func (svc *UserWebService) Update(ctx context.Context, user *model.User) (*model.User, error) {

	// update user
	updUser, err := svc.store.User.Update(ctx, user)
	if err != nil {
		return nil, utils.ErrorWrap(err, "svc.user.Update error")
	}

	return updUser, nil
}

// Delete ...
func (svc *UserWebService) Delete(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	_, err := svc.store.User.Get(ctx, userID)
	if err != nil {
		return utils.ErrorWrap(err, "svc.user.Delete error")
	}

	err = svc.store.User.Delete(ctx, userID)
	if err != nil {
		return utils.ErrorWrap(err, "svc.user.Delete error")
	}

	return nil
}

// List ...
func (svc *UserWebService) List(ctx context.Context) (model.Users, error) {
	users, err := svc.store.User.List(ctx)
	if err != nil {
		return nil, utils.ErrorWrap(err, "svc.user.List")
	}

	return users, nil
}

// GetByEmail ...
func (svc *UserWebService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := svc.store.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, utils.ErrorWrap(err, "svc.user.GetByEmail")
	}

	return user, nil
}