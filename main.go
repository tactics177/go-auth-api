package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/73f2fc18-3053-4c38-943a-416d16432450/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	port := "8080"
	fmt.Println("Server running on port " + port)

	if err := router.Run(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
