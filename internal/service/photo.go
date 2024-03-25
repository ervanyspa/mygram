package service

import (
	"context"
	"mygram/internal/model"
	"mygram/internal/repository"
	"time"
)

type PhotoService interface {
	CreatePhoto(ctx context.Context, photo model.Photo) (model.PhotoCreateRes, error)
	GetPhotosByUserId(ctx context.Context, userId uint64) ([]model.PhotoGetRes, error)
	EditPhoto(ctx context.Context, photo model.Photo) (model.PhotoUpdateRes, error)
	GetPhotoById(ctx context.Context, id uint64) (model.Photo, error)
	DeletePhoto(ctx context.Context, id uint64) error
}

type photoServiceImpl struct {
	repo repository.PhotoQuery
}

func NewPhotoService(repo repository.PhotoQuery) PhotoService {
	return &photoServiceImpl{repo: repo}
}

func (p *photoServiceImpl) CreatePhoto(ctx context.Context, photo model.Photo) (model.PhotoCreateRes, error) {
	res, err := p.repo.CreatePhoto(ctx, photo)
	if err != nil {
		return model.PhotoCreateRes{}, err
	}

	photoResponse := model.PhotoCreateRes{}
	photoResponse.ID = res.ID
	photoResponse.Title = res.Title
	photoResponse.Caption = res.Caption
	photoResponse.PhotoUrl = res.PhotoUrl
	photoResponse.UserId = res.UserId
	photoResponse.CreatedAt = res.CreatedAt

	return photoResponse, nil
}

func (p *photoServiceImpl) GetPhotosByUserId(ctx context.Context, userId uint64) ([]model.PhotoGetRes, error) {
	photos, err := p.repo.GetPhotosByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (p *photoServiceImpl) EditPhoto(ctx context.Context, photo model.Photo) (model.PhotoUpdateRes, error) {
	err := p.repo.EditPhoto(ctx, photo)
	if err != nil {
		return model.PhotoUpdateRes{}, err
	}

	photoResponse := model.PhotoUpdateRes{}
	photoResponse.ID = photo.ID
	photoResponse.Title = photo.Title
	photoResponse.Caption = photo.Caption
	photoResponse.PhotoUrl = photo.PhotoUrl
	photoResponse.UserId = photo.UserId
	photoResponse.UpdatedAt = time.Now()

	return photoResponse, nil
}

func (p *photoServiceImpl) GetPhotoById(ctx context.Context, id uint64) (model.Photo, error) {
	photo, err := p.repo.GetPhotoById(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}

	return photo, nil
}


func (p *photoServiceImpl) DeletePhoto(ctx context.Context, id uint64) error {
	cekPhoto, err := p.repo.GetPhotoById(ctx, id)
	if err != nil {
		return err
	}

	if cekPhoto.ID == 0 {
		return nil
	}
	
	err = p.repo.DeletePhoto(ctx, id)
	if err != nil {
		return err
	}

	return err
}
