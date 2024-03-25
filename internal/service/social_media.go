package service

import (
	"context"
	"mygram/internal/model"
	"mygram/internal/repository"
	"time"
)

type SocialMediaService interface {
	CreateSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMediaCreateRes, error)
	GetSocialMediasByUserId(ctx context.Context, userId uint64) ([]model.SocialMediaGetRes, error)
	EditSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMediaUpdateRes, error)
	GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error)
	DeleteSocialMedia(ctx context.Context, id uint64) error
}

type socialMediaServiceImpl struct {
	repo repository.SocialMediaQuery
}

func NewSocialMediaService(repo repository.SocialMediaQuery) SocialMediaService {
	return &socialMediaServiceImpl{repo: repo}
}

func (s *socialMediaServiceImpl) CreateSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMediaCreateRes, error) {
	res, err := s.repo.CreateSocialMedia(ctx, social)
	if err != nil {
		return model.SocialMediaCreateRes{}, err
	}

	socialMediaResponse := model.SocialMediaCreateRes{}
	socialMediaResponse.ID = res.ID
	socialMediaResponse.Name = res.Name
	socialMediaResponse.SocialMediaUrl = res.SocialMediaUrl
	socialMediaResponse.UserId = res.UserId
	socialMediaResponse.CreatedAt = res.CreatedAt

	return socialMediaResponse, nil
}

func (s *socialMediaServiceImpl) GetSocialMediasByUserId(ctx context.Context, userId uint64) ([]model.SocialMediaGetRes, error) {
	socials, err := s.repo.GetSocialMediasByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return socials, nil
}

func (s *socialMediaServiceImpl) EditSocialMedia(ctx context.Context, social model.SocialMedia) (model.SocialMediaUpdateRes, error) {
	err := s.repo.EditSocialMedia(ctx, social)
	if err != nil {
		return model.SocialMediaUpdateRes{}, err
	}

	socialMediaRes := model.SocialMediaUpdateRes{}
	socialMediaRes.ID = social.ID
	socialMediaRes.Name = social.Name
	socialMediaRes.SocialMediaUrl = social.SocialMediaUrl
	socialMediaRes.UserId = social.UserId
	socialMediaRes.UpdatedAt = time.Now()

	return socialMediaRes, nil
}

func (s *socialMediaServiceImpl) GetSocialMediaById(ctx context.Context, id uint64) (model.SocialMedia, error) {
	social, err := s.repo.GetSocialMediaById(ctx, id)
	if err != nil {
		return model.SocialMedia{}, err
	}

	return social, nil
}



func (s *socialMediaServiceImpl) DeleteSocialMedia(ctx context.Context, id uint64) error {
	cekSocial, err := s.repo.GetSocialMediaById(ctx, id)
	if err != nil {
		return err
	}

	if cekSocial.ID == 0 {
		return nil
	}

	err =  s.repo.DeleteSocialMedia(ctx, id)
	if err != nil {
		return err
	}

	return err
}
