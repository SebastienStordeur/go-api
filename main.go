package main

import (
	controllers "api/controllers/users"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	if err := ConnectToMongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Endpoint working"})
	})
	router.POST("/user", controllers.Signup)

	router.Run("localhost:8080")
}

func ConnectToMongodb() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mongo_uri := os.Getenv("MONGO_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongo_uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	mongoClient = client

	return err
}
