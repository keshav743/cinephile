package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/keshav743/cinephile/routes"
)

func SetupRoutes(app *fiber.App){
	api := app.Group("/api",logger.New())

	routes.SetupUserRoutes(api)
}