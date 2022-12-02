package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/handlers"
	"github.com/keshav743/cinephile/middlewares"
)

func SetupInvitationRoutes(router fiber.Router) {
	friends := router.Group("/invites/friends", middlewares.IsAuthorized)

	friends.Post("/send", middlewares.InvitaionAldreadyExists, handlers.SendFriendRequest)
	friends.Post("/accept", middlewares.InvitaionActionExists, handlers.AcceptFriendRequest)
	friends.Post("/reject", middlewares.InvitaionActionExists, handlers.RejectFriendRequest)
}
