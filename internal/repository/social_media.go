package repository

import (
	"context"
	"mygram/internal/infrastructure"
	"mygram/internal/model"

	"gorm.io/gorm"
)

type SocialMediaQuery interface {
	CreateSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMedia, error)
	GetSocialMediasByUserId(ctx context.Context, userId uint64) ([]model.SocialMediaGetRes, error)
	EditSocialMedia(ctx context.Context, social model.SocialMedia) error
	GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error)
	DeleteSocialMedia(ctx context.Context, id uint64) error
}

type socialMediaQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewSocialMediaQuery(db infrastructure.GormPostgres) SocialMediaQuery {
	return &socialMediaQueryImpl{db: db}
}

func (s *socialMediaQueryImpl) CreateSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMedia, error) {
	db := s.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Create(&social).
		Error; err != nil {
			return model.SocialMedia{}, err
		}

	return social, nil
}

func (s *socialMediaQueryImpl) GetSocialMediasByUserId(ctx context.Context, userId uint64) ([]model.SocialMediaGetRes, error) {
	db := s.db.GetConnection()
	socials := []model.SocialMediaGetRes{}

	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Find(&socials).
		Error; err != nil {
		return nil, err
	}

	return socials, nil
}

func (s *socialMediaQueryImpl) EditSocialMedia(ctx context.Context, social model.SocialMedia) error {
	db := s.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Updates(&social).
		Error; err != nil {
			return err
		}

	return nil
}

func (s *socialMediaQueryImpl) GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error) {
	db := s.db.GetConnection()
	social := model.SocialMedia{}

	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(&social).
		Error; err != nil {
		return model.SocialMedia{}, err
	}

	return social, nil
}

func (s *socialMediaQueryImpl) DeleteSocialMedia(ctx context.Context, id uint64) error {
	db := s.db.GetConnection()

	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Delete(&model.SocialMedia{ID: id}).
		Error; err != nil {
			return err
		}

	return nil
}
