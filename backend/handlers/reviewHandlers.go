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

//TODO - avoid user can review multiple times

func PostReviewMovie(c *fiber.Ctx) error {

	review := new(models.Review)
	movie := new(models.Movie)

	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	handleError(err)

	review.User = id

	err = c.BodyParser(&review)
	handleError(err)

	review.CreatedAt = time.Now()

	err = database.Movies.FindOne(context.TODO(), bson.M{"_id": review.Movie}).Decode(&movie)
	handleError(err)

	result, err := database.Reviews.InsertOne(context.TODO(), review)
	handleError(err)

	newRating := ((movie.Ratings * float64(len(movie.Reviews))) + review.Rating) / float64(len(movie.Reviews)+1)

	database.Movies.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": review.Movie},
		bson.M{"$push": bson.M{"reviews": result.InsertedID}, "$set": bson.M{"ratings": newRating}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id": result.InsertedID,
		},
	})

}
