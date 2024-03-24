package router

import (
	"mygram/internal/handler"
	"mygram/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v *gin.RouterGroup
	handler handler.UserHandler
}

func NewUserRouter(v *gin.RouterGroup, handler handler.UserHandler) UserRouter{
	return &userRouterImpl{v:v, handler: handler}
}

func (u *userRouterImpl) Mount() {
	// activity
	u.v.POST("/register", u.handler.UserSignUp)
	u.v.POST("/login", u.handler.UserSignIn)

	u.v.Use(middleware.CheckAuthBearer)
	u.v.GET("/:id", u.handler.GetUsersById)
	u.v.PUT("/:id", u.handler.EditUser)
	u.v.DELETE("/:id", u.handler.DeleteUsersById)
}

