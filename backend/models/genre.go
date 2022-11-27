package models

type Genre struct {
	ID   int    `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
