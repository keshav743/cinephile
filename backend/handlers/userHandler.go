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
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO - Try to add Cookie-Sessions
//TODO - Handle Duplicate Emails - Signup route
type AuthUser struct {
    Email string `json:"email"`
    Password string `json:"password"`
	Username string `json:"username"`
}

func Signup(c *fiber.Ctx) error {
	user := new(AuthUser)

    if err := c.BodyParser(user); err != nil {
        return err
    }
	hashedPassword, err := helpers.HashPassword(user.Password)
	handleError(err)

	newUser := models.User{
		Username: user.Username,
		Email: user.Email,
		Password: string(hashedPassword),
	}

	database.Users.InsertOne(context.TODO(),newUser)

	token, err := helpers.GenerateJWT(user.Email,user.Username)
	handleError(err)

	c.Response().SetStatusCode(201)
	return c.JSON(
		fiber.Map{
			"status": "success", 
			"data": fiber.Map{
				"username": user.Username,
				"email": user.Email,
				"token": token,
			},
		})
	}

func Signin(c *fiber.Ctx) error {
	user := new(AuthUser)
	existingUser := new(models.User)

    err := c.BodyParser(user)
	handleError(err)

	err = database.Users.FindOne(context.TODO(),bson.M{"email": user.Email}).Decode(&existingUser)
	handleNoDocFoundError(err)

	passwordsMatched := helpers.ComparePasswords(user.Password, existingUser.Password)
	if passwordsMatched != nil {
		c.Response().SetStatusCode(400)
		return c.JSON(
			fiber.Map{
				"status": "failure", 
				"message": "Incorrect credentials",
			})
	}

	token, err := helpers.GenerateJWT(user.Email,user.Username)
	handleError(err)

	c.Response().SetStatusCode(200)
	return c.JSON(
		fiber.Map{
			"status": "success", 
			"data": fiber.Map{
				"username": existingUser.Username,
				"email": existingUser.Email,
				"token": token,
			},
		})
}

func handleNoDocFoundError(err error){
	if err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
}

func handleError(err error){
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(err)
}