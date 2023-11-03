package controllers

import (
	"api/controllers/database"
	"api/models"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Jewel").Collection(collectionName)
	return collection
}

func GetProducts(c *gin.Context) {
	var DB = database.ConnectToDB()
	var collection = GetCollection(DB, "Products")

	/* ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) */

	products, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return
	}
	/* defer cancel() */

	var results models.Product
	for products.Next(context.Background()) {
		err := products.Decode(&results)
		if err != nil {
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"data": results})
}
