package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Username   string               `json:"username" bson:"username,omitempty"`
	Email      string               `json:"email" bson:"email"`
	Password   string               `json:"password,omitempty" bson:"password"`
	Friends    []primitive.ObjectID `json:"friends,omitempty" bson:"friends"`
	Liked      []Movie              `json:"liked,omitempty" bson:"-"`
	Watched    []Movie              `json:"watched,omitempty" bson:"-"`
	Reviews    []Review             `json:"reviews,omitempty" bson:"-"`
	FriendList []User               `json:"friendList,omitempty" bson:"-"`
}
