package routes

import (
	"github.com/tactics177/go-auth-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/73f2fc18-3053-4c38-943a-416d16432450")
	{
		authGroup.POST("/register", controllers.Register)
	}
}
