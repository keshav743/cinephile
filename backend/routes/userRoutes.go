package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/users")

	user.Post("/signup", handlers.Signup)
	user.Post("/signin", handlers.Signin)
}