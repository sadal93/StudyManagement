package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

var assignSurveyCollection = client.Database("test").Collection("assign_surveys")

type AssignedSurvey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	Survey Survey
	UserID string
	StudyID string
	Timestamp int64
	Submitted bool
}

func createAssignSurveyDocument(doc AssignedSurvey) {
	insertResult, err := assignSurveyCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}


