package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// List      primitive.ObjectID  `json:"list" bson:"list"`
	Movie     primitive.ObjectID `json:"movie" bson:"movie"`
	User      primitive.ObjectID `json:"user" bson:"user"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	Review    string             `json:"review" bson:"review"`
	Rating    float64            `json:"rating" bson:"rating"`
}
