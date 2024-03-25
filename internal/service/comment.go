package service

import (
	"context"
	"mygram/internal/model"
	"mygram/internal/repository"
	"time"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment model.Comment) (model.CommentCreateRes, error)
	GetCommentsByPhotoId(ctx context.Context, photoId uint64) ([]model.CommentGetRes, error)
	EditComment(ctx context.Context, comment model.Comment) (model.CommentUpdateRes, error)
	GetCommentById(ctx context.Context, id uint64) (model.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type commentServiceImpl struct {
	repo repository.CommentQuery
}

func NewCommentService(repo repository.CommentQuery) CommentService {
	return &commentServiceImpl{repo: repo}
}

func (c *commentServiceImpl) CreateComment(ctx context.Context, comment model.Comment) (model.CommentCreateRes, error) {
	res, err := c.repo.CreateComment(ctx, comment)
	if err != nil {
		return model.CommentCreateRes{}, err
	}

	commentResponse := model.CommentCreateRes{}
	commentResponse.ID = res.ID
	commentResponse.Message = res.Message
	commentResponse.PhotoId = res.PhotoId
	commentResponse.UserId = res.UserId
	commentResponse.CreatedAt = res.CreatedAt

	return commentResponse, nil
}

func (c *commentServiceImpl) GetCommentsByPhotoId(ctx context.Context, photoId uint64) ([]model.CommentGetRes, error) {
	comments, err := c.repo.GetCommentsByPhotoId(ctx, photoId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *commentServiceImpl) EditComment(ctx context.Context, comment model.Comment) (model.CommentUpdateRes, error) {
	err := c.repo.EditComment(ctx, comment)

	if err != nil {
		return model.CommentUpdateRes{}, err
	}

	commentResponse := model.CommentUpdateRes{}
	commentResponse.ID = comment.ID
	commentResponse.Message = comment.Message
	commentResponse.PhotoId = comment.PhotoId
	commentResponse.UserId = comment.UserId
	commentResponse.UpdatedAt = time.Now()

	return commentResponse, nil
}

func (c *commentServiceImpl) GetCommentById(ctx context.Context, id uint64) (model.Comment, error) {
	comment, err := c.repo.GetCommentById(ctx, id)
	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}



func (c *commentServiceImpl) DeleteComment(ctx context.Context, id uint64) error {
	cekComment, err := c.repo.GetCommentById(ctx, id)
	if err != nil {
		return err
	}

	if cekComment.ID == 0 {
		return nil
	}
	
	err = c.repo.DeleteComment(ctx, id)
	if err != nil {
		return err
	}

	return err
}
