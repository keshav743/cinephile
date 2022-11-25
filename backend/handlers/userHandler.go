package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {
	fmt.Println("User Signup Route")
	err := c.SendString("Hello from signup route")
	return err
}