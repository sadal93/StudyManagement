package main

import (
	pb "StudyManagement/api"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

var assignSurveyCollection = client.Database("test").Collection("assign_surveys")

type AssignedSurvey struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
	survey pb.SurveyData
	UserID string
	StudyID string
}

func createAssignSurveyDocument(doc AssignedSurvey) {
	insertResult, err := assignSurveyCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}


