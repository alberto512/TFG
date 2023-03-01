package jwt

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	SecretKey = []byte(os.Getenv("SECRET_KEY_JWT"))
)

func GenerateToken(username string) (string, error) {
	log.Printf("Generate jwt token")
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Printf("Error: Generating key")
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	log.Printf("Parse jwt token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Printf("Error: Token method not valid")
		   	return nil, errors.New("Unauthorized")
		}
		return SecretKey, nil
	})
	if err != nil {
		log.Printf("Error: Parsing token")
		return "", err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		log.Printf("Error: Token not valid")
		return "", errors.New("Unauthorized")
	}
}