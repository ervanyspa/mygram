package handler

import (
	"mygram/internal/middleware"
	"mygram/internal/model"
	"mygram/internal/service"
	"mygram/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CommentHandler interface {
	CreateComment(ctx *gin.Context)
	GetCommentsByPhotoId(ctx *gin.Context)
	EditComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
}

type commentHandlerImpl struct {
	svc      service.CommentService
	photoSvc service.PhotoService
}

func NewCommentHandler(svc service.CommentService, photoSvc service.PhotoService) CommentHandler {
	return &commentHandlerImpl{svc: svc, photoSvc: photoSvc}
}

func (c *commentHandlerImpl) CreateComment(ctx *gin.Context) {
	userIdClaim, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "cannot claim user ID from token"})
		return
	}
	
	userIdInt, ok := userIdClaim.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Cannot claim user ID from token"})
		return
	}

	commentCreateReq := model.CommentCreateReq{}
	err := ctx.Bind(&commentCreateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(commentCreateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := c.photoSvc.GetPhotoById(ctx, commentCreateReq.PhotoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
		return
	}

	comment := model.Comment{}
	comment.UserId = uint64(userIdInt) 
	comment.Message = commentCreateReq.Message
	comment.PhotoId = commentCreateReq.PhotoId

	commentRes, err := c.svc.CreateComment(ctx, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, commentRes)
}

func (c *commentHandlerImpl) GetCommentsByPhotoId(ctx *gin.Context) {
	photoIdStr := ctx.Request.URL.Query().Get("photo_id")
	if photoIdStr == "" {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "requires photo ID in query"})
		return
	}
	photoId, err := strconv.Atoi(photoIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	comments, err := c.svc.GetCommentsByPhotoId(ctx, uint64(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	
	if len(comments) == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (c *commentHandlerImpl) EditComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if commentId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := c.svc.GetCommentById(ctx, uint64(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment not found"})
		return
	}

	userIdClaim, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user session"})
		return
	}
	
	userIdInt, ok := userIdClaim.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}

	if comment.UserId != uint64(userIdInt)  {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	commentUpdateReq := model.CommentUpdateReq{}
	err = ctx.Bind(&commentUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(commentUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	commentUp := model.Comment{}
	commentUp.ID = uint64(commentId)
	commentUp.UserId = comment.UserId
	commentUp.PhotoId = comment.PhotoId
	commentUp.Message = commentUpdateReq.Message

	commentRes, err := c.svc.EditComment(ctx, commentUp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, commentRes)
}

func (c *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if commentId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := c.svc.GetCommentById(ctx, uint64(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment did not exist"})
		return
	}

	userIdClaim, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	userIdInt, ok := userIdClaim.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
	if comment.UserId != uint64(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	err = c.svc.DeleteComment(ctx, uint64(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "Your comment has been successfully deleted" ,
	})
}