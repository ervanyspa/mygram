package repository

import (
	"context"
	"mygram/internal/infrastructure"
	"mygram/internal/model"

	"gorm.io/gorm"
)

type CommentQuery interface {
	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
	GetCommentsByPhotoId(ctx context.Context, photoId uint64) ([]model.CommentGetRes, error)
	EditComment(ctx context.Context, comment model.Comment) error
	GetCommentById(ctx context.Context, id uint64) (model.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type commentQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewCommentQuery(db infrastructure.GormPostgres) CommentQuery {
	return &commentQueryImpl{db: db}
}

func (c *commentQueryImpl) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
	db := c.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("comments").
		Save(&comment).Error; err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (c *commentQueryImpl) GetCommentsByPhotoId(ctx context.Context, photoId uint64) ([]model.CommentGetRes, error) {
	db := c.db.GetConnection()
	comments := []model.CommentGetRes{}

	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("photo_id = ?", photoId).
		Where("deleted_at IS NULL").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Preload("Photo", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, caption, url, user_id").Table("photos").Where("deleted_at is null")
		}).
		Find(&comments).
		Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *commentQueryImpl) EditComment(ctx context.Context, comment model.Comment) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Updates(&comment).
		Error; err != nil {
			return err
		}

	return nil
}

func (c *commentQueryImpl) GetCommentById(ctx context.Context, id uint64) (model.Comment, error) {
	db := c.db.GetConnection()
	comment := model.Comment{}

	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(&comment).
		Error; err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

func (c *commentQueryImpl) DeleteComment(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("comments").
		Delete(&model.Comment{ID: id}).
		Error;  err != nil {
			return err
		}

	return nil
}
