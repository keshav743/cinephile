package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO - Name from req.params has + when we need to send a req to TMDB

type SearchResponse struct {
	Movies []models.Movie     `json:"results"`
	User   primitive.ObjectID `json:"user"`
	Movie  primitive.ObjectID `json:"movie"`
}

var genrePopulate bson.M = bson.M{"$lookup": bson.M{
	"from":         "Genres",
	"localField":   "genre",
	"foreignField": "_id",
	"as":           "genres",
}}

var watchedPopulate bson.M = bson.M{"$lookup": bson.M{
	"from":         "Users",
	"localField":   "watched",
	"foreignField": "_id",
	"as":           "watched",
}}

var likedPopulate bson.M = bson.M{"$lookup": bson.M{
	"from":         "Users",
	"localField":   "liked",
	"foreignField": "_id",
	"as":           "liked",
}}

var aggProject bson.M = bson.M{"$project": bson.M{
	"title":            1,
	"overview":         1,
	"imageUrl":         1,
	"release":          1,
	"language":         1,
	"popularity":       1,
	"id":               1,
	"tmdbId":           1,
	"genre":            1,
	"watched.username": 1,
	"watched.email":    1,
	"watched._id":      1,
	"liked.email":      1,
	"liked._id":        1,
	"liked.username":   1,
	"reviews":          1,
}}

var TMDBMovieSearchEndpint string = "https://api.themoviedb.org/3/search/movie?api_key=c8af2e4fcd4bf5d99fcb9bfa901fc684"
var wg sync.WaitGroup

func GetMovieByCriteria(c *fiber.Ctx) error {
	var aggSkip, aggLimit bson.M
	movies := make([]bson.M, 10)

	name := c.Query("name")

	url := TMDBMovieSearchEndpint + "&query=" + name + "&page=" + c.Query("page") + "&include_Adult=True"

	result, err := http.Get(url)
	handleError(err)
	defer result.Body.Close()

	var movieResults SearchResponse
	err = json.NewDecoder(result.Body).Decode(&movieResults)
	handleError(err)

	for i := 0; i < len(movieResults.Movies); i++ {
		fmt.Println(movieResults.Movies[i].Title)
		wg.Add(1)
		go func(movie models.Movie) {
			movie.Ratings = 0
			movie.Reviews = make([]primitive.ObjectID, 0)
			movie.Liked = make([]primitive.ObjectID, 0)
			movie.Watched = make([]primitive.ObjectID, 0)
			cnt, err := database.Movies.CountDocuments(context.TODO(), bson.M{"tmdbId": movie.TMDB_ID})
			handleError(err)
			if cnt == 0 {
				database.Movies.InsertOne(context.TODO(), movie)
			}
			wg.Done()
		}(movieResults.Movies[i])
	}

	wg.Wait()

	page, err := strconv.Atoi(c.Query("page"))
	handleError(err)

	aggSearch := bson.M{"$match": bson.M{"title": bson.M{"$regex": name, "$options": "im"}}}

	if page > 0 {
		aggSkip = bson.M{"$skip": int64(int(math.Abs(float64(page-1))) * 20)}
		aggLimit = bson.M{"$limit": 20}
	}

	cursor, err := database.Movies.Aggregate(context.TODO(), []bson.M{
		aggSearch, aggSkip, aggLimit, genrePopulate, aggProject,
	})
	handleError(err)

	err = cursor.All(context.TODO(), &movies)
	handleError(err)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"movies":  movies,
			"message": "movies that were found for your query.",
		},
	})

}

func GetMovieById(c *fiber.Ctx) error {
	movie := make([]bson.M, 1)

	movieId, err := primitive.ObjectIDFromHex(c.Params("id"))
	handleError(err)

	aggSearch := bson.M{"$match": bson.M{"_id": movieId}}

	cursor, _ := database.Movies.Aggregate(context.TODO(), []bson.M{
		aggSearch, genrePopulate, watchedPopulate, likedPopulate, aggProject,
	})

	err = cursor.All(context.TODO(), &movie)
	handleError(err)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"movie":   movie,
			"message": "Requested movie has been found in DB.",
		},
	})
}

func ToggleWatched(c *fiber.Ctx) error {

	toggleWatchedResponse := new(SearchResponse)
	movie := new(models.Movie)

	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	handleError(err)

	err = c.BodyParser(&toggleWatchedResponse)
	handleError(err)

	toggleWatchedResponse.User = id

	err = database.Movies.FindOne(context.TODO(), bson.M{"_id": toggleWatchedResponse.Movie}).Decode(&movie)
	handleError(err)

	for i := 0; i < len(movie.Watched); i++ {
		if movie.Watched[i] == toggleWatchedResponse.User {

			result := database.Movies.FindOneAndUpdate(context.TODO(),
				bson.M{"_id": toggleWatchedResponse.Movie},
				bson.M{"$pull": bson.M{"watched": toggleWatchedResponse.User}},
				options.FindOneAndUpdate().SetReturnDocument(options.After))

			doc := bson.M{}

			decodedErr := result.Decode(&doc)
			handleError(decodedErr)

			return c.JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"movie":   doc,
					"message": "Removed as watched.",
				},
			})
		}
	}

	result := database.Movies.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": toggleWatchedResponse.Movie},
		bson.M{"$push": bson.M{"watched": toggleWatchedResponse.User}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	doc := bson.M{}

	decodedErr := result.Decode(&doc)
	handleError(decodedErr)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"movie":   doc,
			"message": "Marked as watched.",
		},
	})
}

func ToggleLiked(c *fiber.Ctx) error {

	toggleLikedResponse := new(SearchResponse)
	movie := new(models.Movie)

	id, err := primitive.ObjectIDFromHex(c.Locals("id").(string))
	handleError(err)

	err = c.BodyParser(&toggleLikedResponse)
	handleError(err)

	toggleLikedResponse.User = id

	err = database.Movies.FindOne(context.TODO(), bson.M{"_id": toggleLikedResponse.Movie}).Decode(&movie)
	handleError(err)

	for i := 0; i < len(movie.Liked); i++ {
		if movie.Liked[i] == toggleLikedResponse.User {

			result := database.Movies.FindOneAndUpdate(context.TODO(),
				bson.M{"_id": toggleLikedResponse.Movie},
				bson.M{"$pull": bson.M{"liked": toggleLikedResponse.User}},
				options.FindOneAndUpdate().SetReturnDocument(options.After))

			doc := bson.M{}

			decodedErr := result.Decode(&doc)
			handleError(decodedErr)

			return c.JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"movie":   doc,
					"message": "Removed from liked.",
				},
			})
		}
	}

	result := database.Movies.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": toggleLikedResponse.Movie},
		bson.M{"$push": bson.M{"liked": toggleLikedResponse.User}},
		options.FindOneAndUpdate().SetReturnDocument(options.After))

	doc := bson.M{}

	decodedErr := result.Decode(&doc)
	handleError(decodedErr)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"movie":   doc,
			"message": "Marked as liked.",
		},
	})
}
