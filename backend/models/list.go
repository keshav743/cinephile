package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID         primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Access     string               `json:"access" bson:"access"`
	Curator    primitive.ObjectID   `json:"curator" bson:"curator"`
	AccessList []primitive.ObjectID `json:"accessList" bson:"accessList"`
	Movies     []primitive.ObjectID `json:"movies" bson:"movies"`
	Title      string               `json:"title" bson:"title"`
	Overview   string               `json:"overview" bson:"overview"`
	Reviews    []primitive.ObjectID `json:"reviews" bson:"reviews"`
}
