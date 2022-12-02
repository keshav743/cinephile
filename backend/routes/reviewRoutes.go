package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
	"github.com/keshav743/cinephile/middlewares"
)

func SetupReviewRoutes(router fiber.Router) {
	reviews := router.Group("/reviews", middlewares.IsAuthorized)

	reviews.Post("/post", handlers.PostReviewMovie)
}
