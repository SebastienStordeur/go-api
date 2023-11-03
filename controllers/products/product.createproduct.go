package controllers

import (
	"api/controllers/database"
	"api/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(c *gin.Context) {
	var DB = database.ConnectToDB()
	var collection = GetCollection(DB, "Products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	products := new(models.Product)
	defer cancel()

	if err := c.BindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productPayload := models.Product{
		ID:          primitive.NewObjectID(),
		Title:       "A new product",
		Description: "A new description",
		Images:      []string{"image1", "image2"},
		Price:       149,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, productPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	println(result)

	c.JSON(http.StatusCreated, gin.H{"data": result})
	return
}
