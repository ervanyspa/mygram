package router

import (
	"mygram/internal/handler"
	"mygram/internal/middleware"

	"github.com/gin-gonic/gin"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v *gin.RouterGroup
	handler handler.PhotoHandler
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter{
	return &photoRouterImpl{v:v, handler: handler}
}

func (p *photoRouterImpl) Mount() {
	p.v.Use(middleware.CheckAuthBearer)
	p.v.POST("", p.handler.CreatePhoto)
	p.v.GET("",p.handler.GetPhotosByUserId)
	p.v.PUT("/:id",p.handler.EditPhoto)
	p.v.DELETE("/:id", p.handler.DeletePhoto)

}

