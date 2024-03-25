package router

import (
	"mygram/internal/handler"
	"mygram/internal/middleware"

	"github.com/gin-gonic/gin"
)

type CommentRouter interface {
	Mount()
}

type commentRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.CommentHandler
}

func NewCommentRouter(v *gin.RouterGroup, handler handler.CommentHandler) PhotoRouter {
	return &commentRouterImpl{v: v, handler: handler}
}

func (c *commentRouterImpl) Mount() {
	c.v.Use(middleware.CheckAuthBearer)
	c.v.POST("", c.handler.CreateComment)
	c.v.GET("", c.handler.GetCommentsByPhotoId)
	c.v.PUT("/:id", c.handler.EditComment)
	c.v.DELETE("/:id", c.handler.DeleteComment)
}