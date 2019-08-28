package main

import (
	pb "StudyManagement/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var surveyCollection = client.Database("test").Collection("surveys")

type Survey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	Description string
	Type string
	Study string
	Questions[]* pb.Question
}

