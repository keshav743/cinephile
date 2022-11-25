package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(email string, username string) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

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

func IsAuthorized(c *fiber.Ctx) error {

		authToken := c.GetReqHeaders()["Authorization"]
		if authToken == "" {
			c.Response().SetStatusCode(400)
			c.JSON(fiber.Map{
				"status": "failure",
				"message": "Token Missing",
			})
			return nil
		}

		mySigningKey := os.Getenv("JWT_KEY")

		token, err := jwt.ParseWithClaims(authToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})
		HandleError(err)

		if !token.Valid{
			c.Response().SetStatusCode(400)
			c.JSON(fiber.Map{
				"status": "failure",
				"message": "Token Invalid",
			})
			return nil
		}

		return c.Next()
}

func HandleError(err error){
	if err != nil {
		fmt.Println(err)
		return
	}
}