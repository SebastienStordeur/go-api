package controllers

import (
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

	existingUser, err := userCollection.Find(ctx, bson.M{"email": userCredentials.Email})
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user does not exist", "err": err})
		return
	}
	var existingUserStruct models.User
	err = existingUser.Decode(&existingUserStruct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// compare psw
	passwordsMatching := comparePasswords(userCredentials.Password, existingUserStruct.Password)
	if passwordsMatching {
		// The passwords match, so the user is logged in
		c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
	} else {
		// The passwords do not match, so an error is returned
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
	}
}

func comparePasswords(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
