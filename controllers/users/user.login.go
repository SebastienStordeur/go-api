package controllers

import (
	"api/controllers/auth"
	"api/controllers/database"
	"api/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var DB = database.ConnectToDB()
	var userCollection = GetCollection(DB, "Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userCredentials := new(models.User)
	defer cancel()

	if err := c.BindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Search the user in the db
	foundUser, err := userCollection.Find(ctx, bson.M{"email": userCredentials.Email})
	if !foundUser.Next(ctx) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decode the user
	var decodedUser models.User
	err = foundUser.Decode(&decodedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// compare psw
	passwordsMatching := comparePasswords(userCredentials.Password, decodedUser.Password)
	if !passwordsMatching {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to login"})
		return
	}

	// Convert the user ID to a slice of bytes
	userId := decodedUser.ID.Hex()

	// Generate an access token
	tokenChan, err := auth.GenerateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token := <-tokenChan
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func comparePasswords(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
