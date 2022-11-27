package bun

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"gitlab.com/m0ta/lts/app/model"
)

// UserBunRepo ...
type UserBunRepo struct {
	db 		*DB
	logger 	*zap.Logger
}

// NewUserRepo ...
func NewUserRepo(db *DB, logger *zap.Logger) *UserBunRepo {
	return &UserBunRepo{
		db: 	db,
		logger: logger,
	}
}

// Get retrieves user from Postgres
func (repo *UserBunRepo) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	err := repo.db.
		NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		repo.logger.Error("store.User.Get", 
			zap.String("err", err.Error()),
			zap.String("id", id.String()),
		)
		return nil, err
	}
	return user, nil
}

// Create creates user in Postgres
func (repo *UserBunRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.ID = uuid.New()
	_, err := repo.db.
		NewInsert().
		Model(user).
		ExcludeColumn("created_at", "logined_at", "updated_at").
		Returning("*").
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.Create", 
			zap.String("err", err.Error()),
			zap.Object("user", user),
		)
		return nil, err
	}
	return user, nil
}

// Update updates user in Postgres
func (repo *UserBunRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := repo.db.
		NewUpdate().
		Model(user).
		WherePK().
		Returning("*").
		OmitZero().
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.Update", 
			zap.String("err", err.Error()),
			zap.Object("user", user),
		)
		return nil, err
	}

	return user, nil
}

// Delete deletes user from Postgres
func (repo *UserBunRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.
		NewDelete().
		Model((*model.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.Delete", 
			zap.String("err", err.Error()),
			zap.String("id", id.String()),
		)
		return err
	}
	
	return nil
}

// List retrieves user list from Postgres
func (repo *UserBunRepo) List(ctx context.Context) (model.Users, error) {
	var users model.Users
	err := repo.db.
		NewSelect().
		Model(&users).
		Scan(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.List", 
			zap.String("err", err.Error()),
		)
		return nil, err
	}
	return users, nil
}

// GetByEmail retrieves user by email from Postgres
func (repo *UserBunRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := repo.db.
		NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.GetByEmail", 
			zap.String("err", err.Error()),
			zap.String("email", email),
		)
		return nil, err
	}
	return user, nil
}

// Exists checks user if exists from Postgres
func (repo *UserBunRepo) Exists(ctx context.Context, email string) (bool, error) {
	exist, err := repo.db.
		NewSelect().
		//Model(user).
		Model((*model.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	
	if err != nil {
		repo.logger.Error("store.User.Exists", 
			zap.String("err", err.Error()),
			zap.String("email", email),
		)
		return false, err
	}
	return exist, nil
}