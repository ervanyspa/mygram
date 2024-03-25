package router

import (
	"mygram/internal/handler"
	"mygram/internal/middleware"

	"github.com/gin-gonic/gin"
)

type SocialMediaRouter interface {
	Mount()
}

type socialMediaRouterImpl struct {
	v *gin.RouterGroup
	handler handler.SocialMediaHandler
}

func NewSocialMediaRouter(v *gin.RouterGroup, handler handler.SocialMediaHandler) SocialMediaRouter{
	return &socialMediaRouterImpl{v:v, handler: handler}
}

func (s *socialMediaRouterImpl) Mount() {
	s.v.Use(middleware.CheckAuthBearer)
	s.v.POST("", s.handler.CreateSocialMedia)
	s.v.GET("",s.handler.GetSocialMediasByUserId)
	s.v.PUT("/:id",s.handler.EditSocialMedia)
	s.v.DELETE("/:id", s.handler.DeleteSocialMedia)

}

