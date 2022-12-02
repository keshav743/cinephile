package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
	"github.com/keshav743/cinephile/middlewares"
)

func SetupMovieRoutes(router fiber.Router) {
	movies := router.Group("/movies")
	movies.Get("/search", handlers.GetMovieByCriteria)
	movies.Get("/:id", handlers.GetMovieById)

	movies = router.Group("/movies", middlewares.IsAuthorized)
	movies.Post("/status/watched", handlers.ToggleWatched)
	movies.Post("/status/liked", handlers.ToggleLiked)
}
