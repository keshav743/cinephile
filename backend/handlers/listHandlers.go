package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO - Write Middlewares for Checking whether a List exists and whether that List has public access
//TODO - Write Middleware to Check if a movie exists

type ListResponse struct {
	List   primitive.ObjectID `json:"list"`
	Movie  primitive.ObjectID `json:"movie"`
	User   primitive.ObjectID `json:"user"`
	Access string             `json:"access"`
}

func CreateList(c *fiber.Ctx) error {

	list := new(models.List)

	err := c.BodyParser(list)
	handleError(err)

	list.AccessList = make([]primitive.ObjectID, 0)
	list.Reviews = make([]primitive.ObjectID, 0)
	list.Movies = make([]primitive.ObjectID, 0)

	result, _ := database.Lists.InsertOne(context.TODO(), list)
	list.ID = result.InsertedID.(primitive.ObjectID)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"list":    list,
			"message": "List has been created successfully.",
		},
	})
}

func AddMovieToList(c *fiber.Ctx) error {

	addMovieToListResponse := new(ListResponse)

	err := c.BodyParser(addMovieToListResponse)
	handleError(err)

	result := database.Lists.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": addMovieToListResponse.List},
		bson.M{"$push": bson.M{"movies": addMovieToListResponse.Movie}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	doc := bson.M{}

	decodedErr := result.Decode(&doc)
	handleError(decodedErr)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"result":  doc,
			"message": "Movie added to List.",
		},
	})

}

func ChangeAccess(c *fiber.Ctx) error {
	changeAccessResponse := parseBody(c)

	if changeAccessResponse.Access != "public" && changeAccessResponse.Access != "private" {
		return c.JSON(fiber.Map{
			"status":  "failure",
			"message": "Access Parameter has a unknown value.",
		})
	}

	result := database.Lists.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": changeAccessResponse.List},
		bson.M{"$set": bson.M{"access": changeAccessResponse.Access}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	doc := bson.M{}

	decodedErr := result.Decode(&doc)
	handleError(decodedErr)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"list":    doc,
			"message": "Requested list's access has been changed",
		},
	})
}

func GrantAccess(c *fiber.Ctx) error {
	list := new(models.List)
	grantAccessResponse := parseBody(c)

	err := database.Lists.FindOne(context.TODO(), bson.M{"_id": grantAccessResponse.List}).Decode(&list)
	handleNoDocFoundError(err)

	if list.Access == "public" {
		return c.JSON(fiber.Map{
			"status":  "failure",
			"message": "List aldready has public access.",
		})
	}

	for i := 0; i < len(list.AccessList); i++ {
		if list.AccessList[i] == grantAccessResponse.User {
			return c.JSON(fiber.Map{
				"status":  "failure",
				"message": "Request user aldready has access to this list.",
			})
		}
	}

	result := database.Lists.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": grantAccessResponse.List},
		bson.M{"$push": bson.M{"accessList": grantAccessResponse.User}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	doc := bson.M{}

	decodedErr := result.Decode(&doc)
	handleError(decodedErr)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"list":    doc,
			"message": "Requested user has been given access to this list.",
		},
	})
}

func RevokeAccess(c *fiber.Ctx) error {
	list := new(models.List)
	revokeAccessResponse := parseBody(c)

	err := database.Lists.FindOne(context.TODO(), bson.M{"_id": revokeAccessResponse.List}).Decode(&list)
	handleNoDocFoundError(err)

	if list.Access == "public" {
		return c.JSON(fiber.Map{
			"status":  "failure",
			"message": "List has public access. Can't revoke access for a particular user.",
		})
	}

	for i := 0; i < len(list.AccessList); i++ {
		if list.AccessList[i] == revokeAccessResponse.User {

			result := database.Lists.FindOneAndUpdate(context.TODO(),
				bson.M{"_id": revokeAccessResponse.List},
				bson.M{"$pull": bson.M{"accessList": revokeAccessResponse.User}},
				options.FindOneAndUpdate().SetReturnDocument(options.After))

			doc := bson.M{}

			decodedErr := result.Decode(&doc)
			handleError(decodedErr)

			return c.JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"list":    doc,
					"message": "Requested user's access to this list has been revoked.",
				},
			})
		}
	}

	return c.JSON(fiber.Map{
		"status":  "failure",
		"message": "Requested doesn't have access to this list.",
	})

}

func parseBody(c *fiber.Ctx) *ListResponse {
	response := new(ListResponse)
	err := c.BodyParser(&response)
	handleError(err)
	return response
}
