package bun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun/extra/bundebug"

	"gitlab.com/m0ta/lts/app/model"
)

// UserBunRepo ...
type UserBunRepo struct {
	db *DB
}

// NewUserRepo ...
func NewUserRepo(db *DB) *UserBunRepo {
	return &UserBunRepo{db: db}
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
		Exec(ctx)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UserBunRepo updates user in Postgres
func (repo *UserBunRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := repo.db.
		NewUpdate().
		Model(user).
		WherePK().
		//Returning("*").
		OmitZero().
		Exec(ctx)
	
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes user in Postgres
func (repo *UserBunRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.
		NewDelete().
		Model((*model.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	
	if err != nil {
		return err
	}
	
	return nil
}

// List retrieves user from Postgres
func (repo *UserBunRepo) List(ctx context.Context) (model.Users, error) {
	var users model.Users
	err := repo.db.
		NewSelect().
		Model(&users).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByEmail retrieves user from Postgres
func (repo *UserBunRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	bundebug.NewQueryHook(bundebug.WithVerbose(true))
	
	user := &model.User{}
	err := repo.db.
		NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}