package middlewares

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvitationResponse struct {
	InvitationID primitive.ObjectID `json:"invitationId" bson:"invitationId"`
}

func InvitaionAldreadyExists(c *fiber.Ctx) error {

	invitation := new(models.Invitation)

	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	HandleError(err)

	err = c.BodyParser(&invitation)
	HandleError(err)

	invitation.Sender = id

	cnt, err := database.Invitations.CountDocuments(context.TODO(),
		bson.M{"$or": []interface{}{
			bson.M{"sender": invitation.Sender, "receiver": invitation.Receiver},
			bson.M{"sender": invitation.Receiver, "receiver": invitation.Sender},
		}})
	HandleError(err)

	if cnt == 0 {
		return c.Next()
	}

	err = database.Invitations.FindOne(context.TODO(),
		bson.M{"$or": []interface{}{
			bson.M{"sender": invitation.Sender, "receiver": invitation.Receiver},
			bson.M{"sender": invitation.Receiver, "receiver": invitation.Sender},
		}},
		options.FindOne().SetSort(bson.M{"CreatedAt": -1})).Decode(invitation)
	HandleError(err)

	fmt.Println(invitation)

	if invitation.Status == "Declined" {
		fmt.Println("Valid Invitation")
		c.Locals("invitationId", invitation.ID.String())
		return c.Next()
	}

	if invitation.Status == "Accepted" {
		return c.JSON(fiber.Map{
			"status":  "failure",
			"message": "Requested person is already your friend",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "failure",
		"message": "You have a pending request",
	})
}

func InvitaionActionExists(c *fiber.Ctx) error {

	invitationResponse := new(InvitationResponse)
	invitation := new(models.Invitation)

	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	HandleError(err)

	err = c.BodyParser(&invitationResponse)
	HandleError(err)

	err = database.Invitations.FindOne(context.TODO(),
		bson.M{"_id": invitationResponse.InvitationID}).Decode(invitation)
	HandleError(err)

	fmt.Println(invitation)

	if invitation.Receiver != id {
		return c.JSON(fiber.Map{
			"status":  "failure",
			"message": "No such request exists.",
		})
	}

	if invitation.Status == "Pending" {
		fmt.Println("Valid Invitation")
		c.Locals("invitationId", invitation.ID)
		return c.Next()
	}

	return c.JSON(fiber.Map{
		"status":  "failure",
		"message": "You have aldready responded for this request",
	})
}

// TOKENS
//  Keshav - 6388a8a200c5a153d3db28f8 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imtlc2hhdmEwMzAyQGdtYWlsLmNvbSIsImV4cCI6MTY3NTA4NDU5NywiaWQiOiI2Mzg4YThhMjAwYzVhMTUzZDNkYjI4ZjgiLCJ1c2VybmFtZSI6Imtlc2hhdkAyMyJ9.y8SHcZUd9B4LdvB75RaTOiCOHsxqLpFNDIrM_gmLNVg
// 	Krithin - 6388a8dc00c5a153d3db28f9 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImtyaXRoaW4yNkBnbWFpbC5jb20iLCJleHAiOjE2NzUwODQ1NDIsImlkIjoiNjM4OGE4ZGMwMGM1YTE1M2QzZGIyOGY5IiwidXNlcm5hbWUiOiJrcml0aGluQDI2In0.VhExP9mwc768LpbY6R7Bmv40VrYN5tRKKOn4VYySij4
//  Ranga - 6388a87700c5a153d3db28f7 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhbmdhcmFqYW4yMDAyQGdtYWlsLmNvbSIsImV4cCI6MTY3NTA4NDY2MiwiaWQiOiI2Mzg4YTg3NzAwYzVhMTUzZDNkYjI4ZjciLCJ1c2VybmFtZSI6InJhbmdhcmFqYW5AMjAwMiJ9.27bkMz0sur_SlVFocNBQFNdBkZ_qLu2by6ncaEV_Zdw
