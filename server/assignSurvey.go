package main

import (
	"context"
	pb "Thesis-demo/api"
	"log"
)

var assignSurveyCollection = client.Database("test").Collection("assign_surveys")

type AssignSurvey struct {
	AssignmentID int64
	SurveyID int64
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
	assignSurveyDataDoc:= AssignSurvey{assignSurveyData.AssignmentID, assignSurveyData.SurveyID, assignSurveyData.UserID, assignSurveyData.StudyID}
	createAssignSurveyDocument(assignSurveyDataDoc)

	log.Printf("Study Created: %v", assignSurveyData.AssignmentID)
	return &pb.AssignSurveyData{StudyID: assignSurveyData.StudyID, SurveyID: assignSurveyData.SurveyID, UserID: assignSurveyData.UserID, AssignmentID: assignSurveyData.AssignmentID}, nil
}





