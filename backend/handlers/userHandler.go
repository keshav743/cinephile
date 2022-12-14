package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/helpers"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO - Try to add Cookie-Sessions
//TODO - Handle Duplicate Emails and Duplicate Username - Signup route

func Signup(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return err
	}
	hashedPassword, err := helpers.HashPassword(user.Password)
	handleError(err)

	user.Password = string(hashedPassword)
	user.Friends = make([]primitive.ObjectID, 0)

	result, err := database.Users.InsertOne(context.TODO(), user)
	handleError(err)

	token, err := helpers.GenerateJWT(result.InsertedID.(primitive.ObjectID), user.Email, user.Username)
	handleError(err)

	user = helpers.FetchUserDetails(user, result.InsertedID.(primitive.ObjectID))

	c.Response().SetStatusCode(201)

	return c.JSON(
		fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"username": user.Username,
				"email":    user.Email,
				"id":       result.InsertedID,
				"watched":  user.Watched,
				"liked":    user.Liked,
				"reviews":  user.Reviews,
				"friends":  user.FriendList,
				"token":    token,
				"message":  "User signup successfull",
			},
		})
}

func Signin(c *fiber.Ctx) error {
	user := new(models.User)
	existingUser := new(models.User)

	err := c.BodyParser(user)
	handleError(err)

	err = database.Users.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	handleNoDocFoundError(err)

	passwordsMatched := helpers.ComparePasswords(user.Password, existingUser.Password)
	if passwordsMatched != nil {
		c.Response().SetStatusCode(400)
		return c.JSON(
			fiber.Map{
				"status":  "failure",
				"message": "Incorrect credentials",
			})
	}

	token, err := helpers.GenerateJWT(existingUser.ID, user.Email, user.Username)
	handleError(err)

	user.Friends = existingUser.Friends
	user = helpers.FetchUserDetails(user, existingUser.ID)

	c.Response().SetStatusCode(200)
	return c.JSON(
		fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"username": user.Username,
				"email":    user.Email,
				"id":       existingUser.ID,
				"watched":  user.Watched,
				"liked":    user.Liked,
				"reviews":  user.Reviews,
				"friends":  user.FriendList,
				"token":    token,
				"message":  "User signin successfull",
			},
		})
}

func handleNoDocFoundError(err error) {
	if err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(err)
}
