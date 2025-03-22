package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tactics177/go-auth-api/config"
	"github.com/tactics177/go-auth-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)

	router.GET("/73f2fc18-3053-4c38-943a-416d16432450/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is running"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
