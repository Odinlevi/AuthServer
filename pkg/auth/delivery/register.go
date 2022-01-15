package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
)

func RegisterHTTPAuthEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)
}

func RegisterHTTPMessageEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/send", h.send)
	router.GET("/get", h.get)
}
