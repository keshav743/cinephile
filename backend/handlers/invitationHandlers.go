package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO - You cant remove a friend. Fixing that requires changes in middlewares also

func SendFriendRequest(c *fiber.Ctx) error {

	invitation := new(models.Invitation)
	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	handleError(err)

	err = c.BodyParser(&invitation)
	handleError(err)

	invitation.Sender = id
	invitation.Status = "Pending"
	invitation.CreatedAt = time.Now()

	result, err := database.Invitations.InsertOne(context.TODO(), invitation)
	handleError(err)

	invitation.ID = result.InsertedID.(primitive.ObjectID)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"invitation": invitation,
		},
		"message": "Friend request sent.",
	})
}

func AcceptFriendRequest(c *fiber.Ctx) error {

	invitation := new(models.Invitation)

	invitationId := c.Locals("invitationId")

	err := database.Invitations.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": invitationId},
		bson.M{"$set": bson.M{"status": "Accepted"}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(invitation)
	handleError(err)

	database.Users.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": invitation.Sender},
		bson.M{"$push": bson.M{"friends": invitation.Receiver}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	database.Users.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": invitation.Receiver},
		bson.M{"$push": bson.M{"friends": invitation.Sender}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"invitation": invitation,
		},
		"message": "Friend request accepted.",
	})
}

func RejectFriendRequest(c *fiber.Ctx) error {

	invitation := new(models.Invitation)

	invitationId := c.Locals("invitationId")

	err := database.Invitations.FindOneAndUpdate(context.TODO(), bson.M{"_id": invitationId}, bson.M{"$set": bson.M{"status": "Declined"}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(invitation)
	handleError(err)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"invitation": invitation,
		},
		"message": "Friend request declined.",
	})
}
