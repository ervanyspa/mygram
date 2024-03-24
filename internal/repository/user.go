package repository

import (
	"context"
	"mygram/internal/infrastructure"
	"mygram/internal/model"
)

type UserQuery interface {
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUsersByID(ctx context.Context, id uint64) (model.User, error)
	EditUser(ctx context.Context, user *model.User) error
	DeleteUsersByID(ctx context.Context, id uint64) error
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type UserCommand interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type userQueryImpl struct{
	db infrastructure.GromPostgres
}

func NewUserQuery(db infrastructure.GromPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) GetUsersByID(ctx context.Context, id uint64) (model.User, error) {
	db := u.db.GetConnection()
	users := model.User{}
	if err := db.
			WithContext(ctx).
			Table("users").
			Where("id = ?", id).
			Where("deleted_at IS NULL").
			Find(&users).Error; err != nil {
		return model.User{}, nil
	}
	return users, nil
}

func (u *userQueryImpl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	db := u.db.GetConnection()
	user := model.User{}
	if err := db.
			WithContext(ctx).
			Table("users").
			Where("email = ?", email).
			Where("deleted_at IS NULL").
			Find(&user).Error; err != nil {
		return model.User{}, nil
	}
	return user, nil
}

func (u *userQueryImpl) EditUser(ctx context.Context, user *model.User) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Updates(&user).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *userQueryImpl) DeleteUsersByID(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Delete(&model.User{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}