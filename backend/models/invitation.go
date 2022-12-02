package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Sender    primitive.ObjectID `json:"sender" bson:"sender"`
	Receiver  primitive.ObjectID `json:"receiver" bson:"receiver"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
