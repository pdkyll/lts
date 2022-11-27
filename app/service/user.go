package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	//"go.uber.org/zap/zapcore"

	"gitlab.com/m0ta/lts/app/model"
	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

// UserWebService ...
type UserWebService struct {
	ctx   	context.Context
	store 	*store.Store
	logger 	*zap.Logger
}

// NewUserWebService creates a new user web service
func NewUserWebService(ctx context.Context, store *store.Store, logger 	*zap.Logger) *UserWebService {
	return &UserWebService{
		ctx:   	ctx,
		store: 	store,
		logger: logger,
	}
}

// Validate incoming user details...
func (svc UserWebService) Validate(ctx context.Context, user *model.User) error {

	if !strings.Contains(user.Email, "@") {
		return utils.ErrIncorrectEmail
	}

	if len(user.Password) < 6 {
		return utils.ErrIncorrectPassword
	}

	// Email must be unique
	err := svc.Exists(ctx, user.Email)
	if err != nil {
		return utils.ErrUserExists
	}

	return nil
}

// Get ...
func (svc *UserWebService) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := svc.store.User.Get(ctx, id)
	if err != nil {
		return nil, utils.ErrDBError
	}
	
	return user, nil
}

// Create ...
func (svc UserWebService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		svc.logger.Error("utils.HashPassword", 
			zap.String("err", err.Error()),
			zap.String("password", user.Password),
		)
		return nil, utils.ErrNotFound
	}
	user.EncryptedPassword = hash

	user, err = svc.store.User.Create(ctx, user)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return user, nil
}

// Update ...
func (svc *UserWebService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	user, err := svc.store.User.Update(ctx, user)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return user, nil
}

// Delete ...
func (svc *UserWebService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	_, err := svc.Get(ctx, id)
	if err != nil {
		return utils.ErrDBError
	}

	err = svc.store.User.Delete(ctx, id)
	if err != nil {
		return utils.ErrDBError
	}

	return nil
}

// List ...
func (svc *UserWebService) List(ctx context.Context) (model.Users, error) {
	users, err := svc.store.User.List(ctx)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return users, nil
}

// GetByEmail ...
func (svc *UserWebService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := svc.store.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return user, nil
}

// Exists ...
func (svc *UserWebService) Exists(ctx context.Context, email string) error {
	exist, err := svc.store.User.Exists(ctx, email)
	if err != nil {
		return utils.ErrDBError
	}

	if exist {
		return utils.ErrUserExists
	}

	return nil
}

// SignUp ...
func (svc UserWebService) SignUp(ctx context.Context, user *model.User) (*model.User, error) {

	// Create user
	newUser, err := svc.Create(ctx, user)
	if err != nil {
		return nil, utils.ErrDBError
	}

	// Save user token
	token, err := newUser.GenerateToken()//utils.TokenGenerate(newUser)
	if err != nil {
		svc.logger.Error("utils.TokenGenerate", 
			zap.String("err", err.Error()),
			zap.Object("user", newUser),
		)
		return nil, utils.ErrBadToken
	}

	newUser.Token = token

	return newUser, nil
}

// SignIn ...
func (svc UserWebService) SignIn(ctx context.Context, user *model.User) (*model.User, error) {
	// Find user by email
	foundUser, err := svc.GetByEmail(ctx, user.Password)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	// Check password
	err = utils.VerifyPassword(foundUser.EncryptedPassword, user.Password)
	if err != nil {
		svc.logger.Warn("utils.VerifyPassword", 
			zap.String("err", err.Error()),
			zap.Object("user", user),
		)
		return nil, utils.ErrUnmatchPassword
	}

	// Save user token
	token, err := foundUser.GenerateToken()//utils.TokenGenerate(foundUser)
	if err != nil {
		svc.logger.Error("utils.TokenGenerate", 
			zap.String("err", err.Error()),
			zap.Object("user", foundUser),
		)
		return nil, utils.ErrBadToken
	}

	// Log logined time for user
	temp 			:= &model.User{}
	temp.ID 		= user.ID
	temp.LoginedAt 	= time.Now()
	_, err = svc.Update(ctx, temp)
	if err != nil {
		svc.logger.Error("cound't update logined_at", 
			zap.String("err", err.Error()),
			zap.Object("user", temp),
		)
		//return nil, utils.ErrUserNotFound
	}

	foundUser.Token = token

	return foundUser, nil
}

// ChangePassword ...
func (svc UserWebService) ChangePassword(ctx context.Context, user *model.User, password string) (*model.User, error) {
	foundUser, err := svc.Get(ctx, user.ID)
	if err != nil {
		return nil, utils.ErrDBError
	}

	// Check password
	// err = utils.VerifyPassword(foundUser.EncryptedPassword, user.Password)
	// if err != nil {
	// 	svc.logger.Warn("utils.VerifyPassword", 
	// 		zap.String("err", err.Error()),
	// 		zap.Object("user", user),
	// 	)
	// 	return nil, utils.ErrUnmatchPassword
	// }

	hash, err := utils.HashPassword(password)
	if err != nil {
		svc.logger.Error("utils.HashPassword", 
			zap.String("err", err.Error()),
			zap.String("password", password),
		)
		return nil, utils.ErrNotFound

	}
	foundUser.EncryptedPassword = hash

	user, err = svc.Update(ctx, foundUser)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return user, nil
}