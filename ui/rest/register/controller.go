package register

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.POST("/register",handler.CreateUser())
	
}