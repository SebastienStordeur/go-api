package routes

import (
	controllers "api/controllers/users"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/user", controllers.Signup)
	router.POST("/user/login", controllers.Login)
}
