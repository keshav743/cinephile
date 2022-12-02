package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	Title      string               `json:"title" bson:"title"`
	Overview   string               `json:"overview"  bson:"overview"`
	ImageURL   string               `json:"poster_path" bson:"imageUrl"`
	Backdrop   string               `json:"backdrop_path" bson:"backdrop"`
	Release    string               `json:"release_date" bson:"release"`
	Language   string               `json:"original_language" bson:"language"`
	Popularity float64              `json:"popularity" popularity:"popularity"`
	TMDB_ID    int64                `json:"id" bson:"tmdbId"`
	Ratings    float64              `json:"ratings" bson:"ratings"`
	Genre      []int64              `json:"genre_ids" bson:"genre"`
	Watched    []primitive.ObjectID `json:"watched" bson:"watched"`
	Liked      []primitive.ObjectID `json:"liked" bson:"liked"`
	Reviews    []primitive.ObjectID `json:"reviews" bson:"reviews"`
	ID         primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`

	// Genres      []Genre   `json:"genres" bson:"genres"`
	// UserReviews []Reviews `json:"userReviews" bson:"userReviews"`
	// WatchedBy   []Movie   `json:"watchedBy" bson:"watchedBy"`
	// LikedBy     []Movie   `json:"likedBy" bson:"likedBy"`
}
