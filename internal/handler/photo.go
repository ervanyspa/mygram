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

type PhotoHandler interface {
	CreatePhoto(ctx *gin.Context)
	GetPhotosByUserId(ctx *gin.Context)
	EditPhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	svc service.PhotoService
}

func NewPhotoHandler(svc service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{svc: svc}
}

func (p *photoHandlerImpl) CreatePhoto(ctx *gin.Context) {
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

	photoCreateReq := model.PhotoCreateReq{}
	err := ctx.Bind(&photoCreateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(photoCreateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photo := model.Photo{}
	photo.UserId = uint64(userIdInt)
	photo.Caption = photoCreateReq.Caption
	photo.Title = photoCreateReq.Title
	photo.PhotoUrl = photoCreateReq.PhotoUrl

	photoRes, err := p.svc.CreatePhoto(ctx, photo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, photoRes)
}

func (p *photoHandlerImpl) GetPhotosByUserId(ctx *gin.Context) {
	userIdStr := ctx.Request.URL.Query().Get("user_id")
	if userIdStr == "" {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "requires user ID in query"})
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photos, err := p.svc.GetPhotosByUserId(ctx, uint64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(photos) == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

func (p *photoHandlerImpl) EditPhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if photoId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := p.svc.GetPhotoById(ctx, uint64(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
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
	if photo.UserId != uint64(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	photoUpdateReq := model.PhotoUpdateReq{}
	err = ctx.Bind(&photoUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(photoUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photoUp := model.Photo{}
	photoUp.ID = uint64(photoId)
	photoUp.Title = photoUpdateReq.Title
	photoUp.Caption = photoUpdateReq.Caption
	photoUp.PhotoUrl = photoUpdateReq.PhotoUrl
	photoUp.UserId = photo.UserId

	photoRes, err := p.svc.EditPhoto(ctx, photoUp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photoRes)
}

func (p *photoHandlerImpl) DeletePhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if photoId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := p.svc.GetPhotoById(ctx, uint64(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
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
	if photo.UserId != uint64(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	err = p.svc.DeletePhoto(ctx, uint64(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "Your photo has been successfully deleted" ,
	})
}
