package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	Access     string               `json:"access" bson:"access"`
	Title      string               `json:"title" bson:"title"`
	Overview   string               `json:"overview" bson:"overview"`
	CreatedAt  time.Time            `json:"createdAt" bson:"createdAt"`
	AccessList []primitive.ObjectID `json:"accessList" bson:"accessList"`
	Movies     []primitive.ObjectID `json:"movies" bson:"movies"`
	Reviews    []primitive.ObjectID `json:"reviews" bson:"reviews"`
	Curator    primitive.ObjectID   `json:"curator" bson:"curator"`
	ID         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
}
