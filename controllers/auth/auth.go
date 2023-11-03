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

func GenerateAccessToken(id string) (string, error) {
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
		return "", err
	}

	return signedToken, nil
}

func DecodeJwtToken(tokenString string, verificationKey []byte) (*jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return verificationKey, nil
	})
	if err != nil {
		return nil, err
	}
	println(token)

	return &claims, nil
}
