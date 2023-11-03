package auth

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtToken struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaim struct {
	ID string
	jwt.StandardClaims
}

func GenerateAccessToken(id string) (chan string, error) {

	tokenChan := make(chan string)

	go func() {
		claims := *&JwtClaim{
			ID: id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
				Issuer:    "JWT",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Println(err)
			return
		}

		tokenChan <- signedToken
		close(tokenChan)
	}()
	return tokenChan, nil
}
