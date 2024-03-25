package main

import (
	"fmt"
	"mygram/internal/handler"
	"mygram/internal/infrastructure"
	"mygram/internal/model"
	"mygram/internal/repository"
	"mygram/internal/router"
	"mygram/internal/service"
	"mygram/pkg"
	"mygram/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			GO PROJECT STRUCTURE
// @version		2.0
// @description	golang project structure
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
// @schemes		http
func main() {
	g := gin.Default()
	// requirement technical:
	// [x] middleware untuk recover ketika panic
	// [x] mengecheck basic auth

	g.Use(gin.Recovery())

	// /public => generate JWT public
	g.GET("/public", func(ctx *gin.Context) {
		now := time.Now()

		claim := model.StandardClaim{
			Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
			Iss: "go-middleware",
			Aud: "golang-006",
			Sub: "public-token",
			Exp: uint64(now.Add(time.Hour).Unix()),
			Iat: uint64(now.Unix()),
			Nbf: uint64(now.Unix()),
		}
		token, err := helper.GenerateToken(claim)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Message: "error generating public token",
				Errors:  []string{err.Error()},
			})
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{"token": token})
	})

	gorm := infrastructure.NewGormPostgres()
	
	// usersGroup.Use(middleware.CheckAuthBasic)
	// usersGroup.Use(middleware.CheckAuthBearer)
	
	
	usersGroup := g.Group("/users")
	// dependency injection
	userRepo := repository.NewUserQuery(gorm)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)
	userRouter := router.NewUserRouter(usersGroup, userHdl)
	// mount
	userRouter.Mount()


	photosGroup := g.Group("/photos")
	photoRepo := repository.NewPhotoQuery(gorm)
	photoSvc := service.NewPhotoService(photoRepo)
	photoHdl := handler.NewPhotoHandler(photoSvc)
	photoRouter := router.NewPhotoRouter(photosGroup, photoHdl)
	photoRouter.Mount()

	commentsGroup := g.Group("/comments")
	commentRepo := repository.NewCommentQuery(gorm)
	commentSvc := service.NewCommentService(commentRepo)
	commentHdl := handler.NewCommentHandler(commentSvc, photoSvc)
	commentRouter := router.NewCommentRouter(commentsGroup, commentHdl)
	commentRouter.Mount()

	

	
	// swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Run(":3000")
	// Product:
	// authorization menggunakan jwt
	// authentication bisa dilakukan dengan login
	// ketika user login, akan memunculkan JWT ketika success
}