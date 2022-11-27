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
)

//TODO - Name from req.params has + when we need to send a req to TMDB

type SearchResponse struct {
	Movies []models.Movie `json:"results"`
}

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

	genrePopulate := bson.M{"$lookup": bson.M{
		"from":         "Genres",
		"localField":   "genre",
		"foreignField": "_id",
		"as":           "genres",
	}}

	aggProject := bson.M{"$project": bson.M{
		"title":      1,
		"overview":   1,
		"imageUrl":   1,
		"release":    1,
		"language":   1,
		"popularity": 1,
		"id":         1,
		"tmdb_id":    1,
		"genres":     1,
	}}

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
		aggSearch,
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
