package main

import (
	controllers "api/controllers/users"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user", controllers.Signup)
	router.POST("/user/login", controllers.Login)

	router.Run("localhost:8080")
}
