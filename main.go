package main

import (
	middlewares "api/middlewares/cors"
	"api/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	router.GET("/", func (c *gin.Context)  {
		c.JSON(http.StatusAccepted, gin.H{"message": "ALL WORKING"})
		
	})

	routes.UserRoutes(router)
	routes.ProductRoutes(router)

	/* router.Use(auth.AuthMiddleware()) */
	router.GET("/private", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		c.JSON(http.StatusOK, gin.H{"message": "This is a private route", "userID": userID})
	})

	router.Run("localhost:8080")
}
