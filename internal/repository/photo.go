package repository

import (
	"context"
	"mygram/internal/infrastructure"
	"mygram/internal/model"

	"gorm.io/gorm"
)

type PhotoQuery interface {
	CreatePhoto(ctx context.Context, photo model.Photo) (model.Photo, error)
	GetPhotosByUserId(ctx context.Context, userId uint64) ([]model.PhotoGetRes, error)
	EditPhoto(ctx context.Context, photo model.Photo) error
	GetPhotoById(ctx context.Context, id uint64) (model.Photo, error)
	DeletePhoto(ctx context.Context, id uint64) error
}

type photoQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewPhotoQuery(db infrastructure.GormPostgres) PhotoQuery {
	return &photoQueryImpl{db: db}
}

func (p *photoQueryImpl) CreatePhoto(ctx context.Context, photo model.Photo) (model.Photo, error) {
	db := p.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("photos").
		Save(&photo).Error; err != nil {
		return model.Photo{},err
	}
	return photo,nil
}

func (p *photoQueryImpl) GetPhotosByUserId(ctx context.Context, userId uint64) ([]model.PhotoGetRes, error) {
	db := p.db.GetConnection()
	photos := []model.PhotoGetRes{}

	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Find(&photos).
		Error; err != nil {
		return nil, err
	}

	return photos, nil
}

func (p *photoQueryImpl) EditPhoto(ctx context.Context, photo model.Photo) error {
	db := p.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Updates(&photo).
		Error; err != nil {
			return err
		}
	return nil
}

func (p *photoQueryImpl) GetPhotoById(ctx context.Context, id uint64) (model.Photo, error) {
	db := p.db.GetConnection()
	photo := model.Photo{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(&photo).
		Error; err != nil {
		return model.Photo{}, err
	}

	return photo, nil
}



func (p *photoQueryImpl) DeletePhoto(ctx context.Context, id uint64) error {
	db := p.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("photos").
		Delete(&model.Photo{ID: id}).
		Error;  err != nil {
			return err
		}

	return nil
}
