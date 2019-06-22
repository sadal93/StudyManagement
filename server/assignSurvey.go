package main

import (
	"context"
	pb "Thesis-demo/api"
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

func (s *server)  AssignSurvey(ctx context.Context, assignSurveyData *pb.AssignSurveyData) (*pb.AssignSurveyData, error) {
	assignSurveyDataDoc:= AssignSurvey{primitive.NewObjectID(), assignSurveyData.SurveyID, assignSurveyData.UserID, assignSurveyData.StudyID}
	createAssignSurveyDocument(assignSurveyDataDoc)

	log.Printf("Assignment Created: %v", assignSurveyData.Id)
	return &pb.AssignSurveyData{StudyID: assignSurveyData.StudyID, SurveyID: assignSurveyData.SurveyID, UserID: assignSurveyData.UserID, Id: assignSurveyData.Id}, nil
}





