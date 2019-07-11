package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

var assignSurveyCollection = client.Database("test").Collection("assign_surveys")

type AssignSurvey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	SurveyID string
	UserID string
	StudyID string
}

func createAssignSurveyDocument(doc AssignSurvey) {
	insertResult, err := assignSurveyCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}


