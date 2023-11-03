package controllers

import (
	"api/controllers/database"
	"api/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Jewel").Collection(collectionName)
	return collection
}

func Signup(c *gin.Context) {
	var DB = database.ConnectToDB()
	var collection = GetCollection(DB, "Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(models.User)
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	userExists, err := collection.Find(ctx, bson.M{"email": user.Email})
	println(err)

	// If the user does not exist, create them
	if !userExists.Next(ctx) {
		// Hash the password before being stored in the database
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured. Try again later"})
			return
		}

		userPayload := models.User{
			ID:        primitive.NewObjectID(),
			Email:     user.Email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := collection.InsertOne(ctx, userPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		println(result)

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
		return
	}

	// If the user already exists, return an error
	c.JSON(http.StatusConflict, gin.H{"error": "This user already exists"})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return password, err
	}
	return string(hashedPassword), nil
}
