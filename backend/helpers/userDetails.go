package helpers

import (
	"context"
	"sync"

	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchUserDetails(user *models.User, id primitive.ObjectID) *models.User {

	var watched []models.Movie = make([]models.Movie, 0)
	var liked []models.Movie = make([]models.Movie, 0)
	var reviews []models.Review = make([]models.Review, 0)
	var friends []models.User = make([]models.User, 0)

	var findQuery []primitive.ObjectID = append(make([]primitive.ObjectID, 0), id)

	var watchedSearch bson.M = bson.M{"$match": bson.M{"watched": bson.M{"$all": findQuery}}}
	var likedSearch bson.M = bson.M{"$match": bson.M{"liked": bson.M{"$all": findQuery}}}
	var reviewSearch bson.M = bson.M{"$match": bson.M{"user": bson.M{"$all": findQuery}}}

	var userProject bson.M = bson.M{"$project": bson.M{
		"_id":      1,
		"username": 1,
		"email":    1,
	}}

	var aggMovieProject bson.M = bson.M{"$project": bson.M{
		"title":      1,
		"overview":   1,
		"imageUrl":   1,
		"release":    1,
		"language":   1,
		"popularity": 1,
		"_id":        1,
		"tmdbId":     1,
		"genre":      1,
		"watched":    1,
		"liked":      1,
		"reviews":    1,
		"ratings":    1,
	}}

	var wg sync.WaitGroup

	wg.Add(4)

	go func(friendsList []primitive.ObjectID) {

		cursor, err := database.Users.Aggregate(context.TODO(), []bson.M{
			{"$match": bson.M{"_id": bson.M{"$in": friendsList}}},
			userProject,
		})
		HandleError(err)

		err = cursor.All(context.TODO(), &friends)
		HandleError(err)
		wg.Done()
	}(user.Friends)

	go func(id primitive.ObjectID) {
		cursor, err := database.Movies.Aggregate(context.TODO(), []bson.M{
			watchedSearch, aggMovieProject,
		})
		HandleError(err)

		err = cursor.All(context.TODO(), &watched)
		HandleError(err)
		wg.Done()
	}(id)

	go func(id primitive.ObjectID) {
		cursor, err := database.Movies.Aggregate(context.TODO(), []bson.M{
			likedSearch, aggMovieProject,
		})
		HandleError(err)

		err = cursor.All(context.TODO(), &liked)
		HandleError(err)
		wg.Done()
	}(id)

	go func(id primitive.ObjectID) {
		cursor, err := database.Reviews.Aggregate(context.TODO(), []bson.M{
			reviewSearch,
		})
		HandleError(err)

		err = cursor.All(context.TODO(), &reviews)
		HandleError(err)
		wg.Done()
	}(id)

	wg.Wait()

	user.Watched = watched
	user.Liked = liked
	user.Reviews = reviews
	user.FriendList = friends

	return user
}
