package routes

import (
	controllers "api/controllers/products"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	router.GET("/products", controllers.GetProducts)
	router.POST("/products", controllers.CreateProduct)
}
