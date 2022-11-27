package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
)

func SetupMovieRoutes(router fiber.Router) {
	movies := router.Group("/movies")

	movies.Get("/search", handlers.GetMovieByCriteria)
	movies.Get("/:id", handlers.GetMovieById)
}
