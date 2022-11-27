package service

import (
	"context"
	"time"
	//"strings"

	"go.uber.org/zap"
	
	"gitlab.com/m0ta/lts/app/model"
	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

// TokenWebService ...
type TokenWebService struct {
	ctx   	context.Context
	store 	*store.Store
	logger 	*zap.Logger
}

// NewTokenWebService creates a new token web service
func NewTokenWebService(ctx context.Context, store *store.Store, logger *zap.Logger) *TokenWebService {
	return &TokenWebService{
		ctx:   	ctx,
		store: 	store,
		logger: logger,
	}
}

// Validate ...
func (svc *TokenWebService) Validate(ctx context.Context, token *model.Token) error {
	// TODO Checks for correct domain
	// if !strings.Contains(token.Domain, ".") {
	// 	return utils.ErrIncorrectDomain
	// }
	if !(len(token.Domain) > 0) {
		return utils.ErrIncorrectDomain
	}

	return nil
}

// Get ...
func (svc *TokenWebService) Get(ctx context.Context, id uint64) (*model.Token, error) {
	token, err := svc.store.Token.Get(ctx, id)
	if err != nil {
		return nil, utils.ErrDBError
	}

	// Log request to token
	temp 			:= &model.Token{}
	temp.ID 		= id
	temp.RequestedAt= time.Now()
	_, err = svc.Update(ctx, temp)
	if err != nil {
		svc.logger.Error("cound't update requested_at", 
			zap.String("err", err.Error()),
			zap.Object("token", temp),
		)
		//return nil, err
	} else {
		token.RequestedAt = temp.RequestedAt
	}
	
	return token, nil
}

// Create ...
func (svc *TokenWebService) Create(ctx context.Context, token *model.Token) (*model.Token, error) {
	token.CreatedAt = time.Now()
	err := token.GenerateToken()
	if err != nil {
		svc.logger.Error("cound't generate token", 
			zap.String("err", err.Error()),
			zap.Object("token", token),
		)
		return nil, utils.ErrBadToken
	}
	token, err = svc.store.Token.Create(ctx, token)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return token, nil
}

// Update ...
func (svc *TokenWebService) Update(ctx context.Context, token *model.Token) (*model.Token, error) {
	token, err := svc.store.Token.Update(ctx, token)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return token, nil
}

// Delete ...
func (svc *TokenWebService) Delete(ctx context.Context, id uint64) error {
	// Check if token exists
	_, err := svc.store.Token.Get(ctx, id)
	if err != nil {
		return utils.ErrDBError
	}

	err = svc.store.Token.Delete(ctx, id)
	if err != nil {
		return utils.ErrDBError
	}

	return nil
}

// List ...
func (svc *TokenWebService) List(ctx context.Context) (model.Tokens, error) {
	tokens, err := svc.store.Token.List(ctx)
	if err != nil {
		return nil, utils.ErrDBError
	}

	return tokens, nil
}

// GetByDomain ...
func (svc *TokenWebService) GetByDomain(ctx context.Context, domain string) (*model.Token, error) {
	token, err := svc.store.Token.GetByDomain(ctx, domain)
	if err != nil {
		return nil, utils.ErrDBError
	}

	// Log request to token
	temp 			:= &model.Token{}
	temp.ID 		= token.ID
	temp.RequestedAt= time.Now()
	_, err = svc.Update(ctx, temp)
	if err != nil {
		svc.logger.Error("cound't update requested_at", 
			zap.String("err", err.Error()),
			zap.Object("token", temp),
		)
		//return nil, err
	} else {
		token.RequestedAt = temp.RequestedAt
	}

	return token, nil
}

// Exists ...
func (svc *TokenWebService) Exists(ctx context.Context, domain string) error {
	exist, err := svc.store.Token.Exists(ctx, domain)
	if err != nil {
		return utils.ErrDBError
	}
	if !exist {
		return utils.ErrBadToken
	}

	return nil
}

// Verify ...
func (svc *TokenWebService) Verify(ctx context.Context, token *model.Token) (*model.ValidData, error) {
	data := &model.ValidData{}
	
	if (token.ExpiredAt.After(time.Now())) {
		data.Valid = true
	} else {
		data.Valid = false
	}

	data.Token = token
	return data, nil
}