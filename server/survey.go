package main

import (
	pb "StudyManagement/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var surveyCollection = client.Database("test").Collection("survey")

type Survey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	Description string
	Questions[]* pb.Question
}

