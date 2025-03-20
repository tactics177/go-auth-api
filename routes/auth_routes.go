package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tactics177/go-auth-api/internal/handlers"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/73f2fc18-3053-4c38-943a-416d16432450")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/forgot-password", handlers.ForgotPassword)
		authGroup.POST("/reset-password", handlers.ResetPassword)
	}
}
