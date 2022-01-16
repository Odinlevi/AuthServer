package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
)

func RegisterHTTPAuthEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)

	router.GET("/confirm/:token", h.confirm)
}

func RegisterHTTPMessageEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/send", h.send)
	router.GET("/get", h.get)
}

func RegisterHTTPChatEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/add_companion/:companion", h.addCompanion)
	router.POST("/remove_companion/:companion", h.removeCompanion)
	router.GET("/get_companions", h.getCompanions)
}
