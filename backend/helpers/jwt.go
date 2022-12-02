package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateJWT(id primitive.ObjectID, email string, username string) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 86400).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Token creation failed: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
