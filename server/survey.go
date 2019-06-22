package main

import "go.mongodb.org/mongo-driver/bson/primitive"

var surveyCollection = client.Database("test").Collection("users")

type Survey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	Type string
	Study string
	Description string
}

