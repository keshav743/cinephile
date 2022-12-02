package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsAuthorized(c *fiber.Ctx) error {

	authToken := c.GetReqHeaders()["Authorization"]
	if authToken == "" {
		c.Response().SetStatusCode(400)
		c.JSON(fiber.Map{
			"status":  "failure",
			"message": "Token Missing",
		})
		return nil
	}

	claims := jwt.MapClaims{}
	mySigningKey := os.Getenv("JWT_KEY")

	token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	HandleError(err)

	if !token.Valid {
		c.Response().SetStatusCode(400)
		c.JSON(fiber.Map{
			"status":  "failure",
			"message": "Token Invalid",
		})
		return nil
	}

	for key, val := range claims {
		if key == "id" {
			c.Locals(key, val.(string))
		}
	}

	return c.Next()
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
