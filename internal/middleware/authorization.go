package middleware

import (
	"encoding/base64"
	"fmt"
	"mygram/pkg"
	"mygram/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	STATIC_USERNAME = "golang006awesome"
	STATIC_PASSWORD = "mysecretpassword"

	CLAIM_USER_ID  = "claim_user_id"
	CLAIM_USERNAME = "claim_username"
)

func CheckAuthBasic(ctx *gin.Context) {
	// check authorization request
	// step1: ambil data auth dari header
	auth := ctx.GetHeader("Authorization")
	// "Basic Z29sYW5nMDA2YXdlc29tZTpteXNlY3JldHBhc3N3b3Jk"
	// step2: dapatkan base64 string
	authArr := strings.Split(auth, " ")
	if len(authArr) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token"},
		})
		return
	}
	if authArr[0] != "Basic" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid authorization method"},
		})
		return
	}
	// step3: decode base64 string
	token := authArr[1]
	basic, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token", "failed to decode"},
		})
		return
	}
	// step4: compare dengan variable static
	if string(basic) != fmt.Sprintf("%v:%v", STATIC_USERNAME, STATIC_PASSWORD) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid username or password"},
		})
		return
	}
	ctx.Next()
}

func CheckAuthBearer(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")

	authArr := strings.Split(auth, " ")
	if len(authArr) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token"},
		})
		return
	}
	if authArr[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid authorization method"},
		})
		return
	}

	token := authArr[1]
	claims, err := helper.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, pkg.ErrorResponse{
			Message: "unauthorized",
			Errors:  []string{"invalid token", "failed to decode"},
		})
		return
	}
	ctx.Set(CLAIM_USER_ID, claims["user_id"])
	ctx.Set(CLAIM_USERNAME, claims["username"])
	ctx.Next()
}