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
	"go.mongodb.org/mongo-driver/mongo"
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

	// Find user in the db
	user, err := findUserByEmail(ctx, userCollection, userCredentials.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// compare psw
	passwordsMatching := comparePasswords(userCredentials.Password, user.Password)
	if !passwordsMatching {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to login"})
		return
	}

	// Generate an access token
	userId := user.ID.Hex()
	tokenChan := make(chan string)
	go func() {
		token, err := auth.GenerateAccessToken(userId)
		if err != nil {
			tokenChan <- ""
			return
		}
		tokenChan <- token
	}()
	token := <-tokenChan

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func findUserByEmail(ctx context.Context, collection *mongo.Collection, email string) (*models.User, error) {
	user := &models.User{}
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func comparePasswords(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
