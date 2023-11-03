package main

import (
	auth "api/middlewares"
	"api/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.UserRoutes(router)

	router.Use(auth.AuthMiddleware())
	router.GET("/private", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		c.JSON(http.StatusOK, gin.H{"message": "This is a private route", "userID": userID})
	})

	router.Run("localhost:8080")
}
