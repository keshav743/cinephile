package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
	"github.com/keshav743/cinephile/middlewares"
)

func SetupListRoutes(router fiber.Router) {
	list := router.Group("/list", middlewares.IsAuthorized)

	list.Post("/create", handlers.CreateList)
	list.Post("/add", handlers.AddMovieToList)
	list.Post("/access/change", handlers.ChangeAccess)
	list.Post("/access/grant", handlers.GrantAccess)
	list.Post("/access/revoke", handlers.RevokeAccess)
}
