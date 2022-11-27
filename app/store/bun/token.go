package bun

import (
	"context"
	
	"go.uber.org/zap"

	"gitlab.com/m0ta/lts/app/model"
)

// TokenBunRepo ...
type TokenBunRepo struct {
	db *DB
	logger *zap.Logger
}

// NewTokenRepo ...
func NewTokenRepo(db *DB, logger *zap.Logger) *TokenBunRepo {
	return &TokenBunRepo{
		db: 	db,
		logger: logger,
	}
}

// Get retrieves token from Postgres
func (repo *TokenBunRepo) Get(ctx context.Context, id uint64) (*model.Token, error) {
	token := &model.Token{}
	err := repo.db.
		NewSelect().
		Model(token).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		repo.logger.Error("store.Token.Get", 
			zap.String("err", err.Error()),
			zap.Uint64("id", id),
		)
		return nil, err
	}

	return token, nil
}

// Create creates token in Postgres
func (repo *TokenBunRepo) Create(ctx context.Context, token *model.Token) (*model.Token, error) {
	_, err := repo.db.
		NewInsert().
		Model(token).
		ExcludeColumn("id", "updated_at", "requested_at").
		Returning("*").
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.Create", 
			zap.String("err", err.Error()),
			zap.Object("token", token),
		)
		return nil, err
	}
	return token, nil
}

// Update updates token in Postgres
func (repo *TokenBunRepo) Update(ctx context.Context, token *model.Token) (*model.Token, error) {
	_, err := repo.db.
		NewUpdate().
		Model(token).
		WherePK().
		//Returning("*").
		OmitZero().
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.Update", 
			zap.String("err", err.Error()),
			zap.Object("token", token),
		)
		return nil, err
	}

	return token, nil
}

// Delete deletes token from Postgres
func (repo *TokenBunRepo) Delete(ctx context.Context, id uint64) error {
	_, err := repo.db.
		NewDelete().
		Model((*model.Token)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.Delete", 
			zap.String("err", err.Error()),
			zap.Uint64("id", id),
		)
		return err
	}
	
	return nil
}

// List retrieves token list from Postgres
func (repo *TokenBunRepo) List(ctx context.Context) (model.Tokens, error) {
	var tokens model.Tokens
	err := repo.db.
		NewSelect().
		Model(&tokens).
		Scan(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.List", 
			zap.String("err", err.Error()),
		)
		return nil, err
	}
	return tokens, nil
}

// GetByDomain retrieves token by domain from Postgres
func (repo *TokenBunRepo) GetByDomain(ctx context.Context, domain string) (*model.Token, error) {
	token := &model.Token{}
	err := repo.db.
		NewSelect().
		Model(token).
		Where("domain = ?", domain).
		Scan(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.GetByDomain", 
			zap.String("err", err.Error()),
			zap.String("domain", domain),
		)
		return nil, err
	}
	return token, nil
}

// Exists checks token if exists from Postgres
func (repo *TokenBunRepo) Exists(ctx context.Context, domain string) (bool, error) {
	exist, err := repo.db.
		NewSelect().
		Model((*model.Token)(nil)).
		Where("domain = ?", domain).
		Exists(ctx)
	
	if err != nil {
		repo.logger.Error("store.Token.Exists", 
			zap.String("err", err.Error()),
			zap.String("domain", domain),
		)
		return false, err
	}
	return exist, nil
}