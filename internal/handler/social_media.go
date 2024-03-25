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

type SocialMediaHandler interface {
	CreateSocialMedia(ctx *gin.Context)
	GetSocialMediasByUserId(ctx *gin.Context)
	EditSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
	svc service.SocialMediaService
}

func NewSocialMediaHandler(svc service.SocialMediaService) SocialMediaHandler {
	return &socialMediaHandlerImpl{svc: svc}
}

func (s *socialMediaHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
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

	socialMediaReq := model.SocialMediaReq{}
	err := ctx.Bind(&socialMediaReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(socialMediaReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	socialmedia := model.SocialMedia{}
	socialmedia.Name = socialMediaReq.Name
	socialmedia.SocialMediaUrl = socialMediaReq.SocialMediaUrl
	socialmedia.UserId = uint64(userIdInt)

	socialMediaRes, err := s.svc.CreateSocialMedia(ctx, socialmedia)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, socialMediaRes)
}

func (s *socialMediaHandlerImpl) GetSocialMediasByUserId(ctx *gin.Context) {
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

	socials, err := s.svc.GetSocialMediasByUserId(ctx, uint64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(socials) == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "social media not found"})
		return
	}

	ctx.JSON(http.StatusOK, socials)
}

func (s *socialMediaHandlerImpl) EditSocialMedia(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if socialId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	social, err := s.svc.GetSocialMediaById(ctx, uint64(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if social.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "social media not found"})
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
	if social.UserId != uint64(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	socialMediaUpdateReq := model.SocialMediaReq{}
	err = ctx.Bind(&socialMediaUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(socialMediaUpdateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	socialmediaUp := model.SocialMedia{}
	socialmediaUp.ID = uint64(socialId)
	socialmediaUp.Name = socialMediaUpdateReq.Name
	socialmediaUp.SocialMediaUrl = socialMediaUpdateReq.SocialMediaUrl
	socialmediaUp.UserId = social.UserId

	socialMediaRes, err := s.svc.EditSocialMedia(ctx, socialmediaUp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, socialMediaRes)
}

func (s *socialMediaHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if socialId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	social, err := s.svc.GetSocialMediaById(ctx, uint64(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	if social.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social Media not found"})
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
	if social.UserId != uint64(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	err = s.svc.DeleteSocialMedia(ctx, uint64(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "Your social media has been successfully deleted" ,
	})
}
