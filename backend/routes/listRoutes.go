package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
)

func SetupListRoutes(router fiber.Router) {
	list := router.Group("/list")

	list.Post("/create", handlers.CreateList)
	list.Post("/add", handlers.AddMovieToList)
	list.Post("/access/change", handlers.ChangeAccess)
	list.Post("/access/grant", handlers.GrantAccess)
	list.Post("/access/revoke", handlers.RevokeAccess)
}
