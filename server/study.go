package main

import (
	pb "Thesis-demo/api"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

var studyCollection = client.Database("test").Collection("study")

type Study struct {
	StudyID string
	Name string
	Description string
	StartDate int64
	Status string
	Users []string
}

func createStudyDocument(doc Study) {
	insertResult, err := studyCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}

func (s *server)  CreateStudy(ctx context.Context, study *pb.StudyMetaData) (*pb.StudyMetaData, error) {

	studyDoc:= Study{study.StudyID, study.Name, study.Description, study.StartDate, study.Status,study.Users}
	createStudyDocument(studyDoc)

	log.Printf("Study Created: %v", study.Name)
	return &pb.StudyMetaData{StudyID: study.StudyID, Name: study.Name, Description: study.Description, StartDate: study.StartDate, Status: study.Status, Users: study.Users}, nil
}

func (s *server)  GetUsers(ctx context.Context, study *pb.StudyMetaData) (*pb.StudyUsers, error) {

	var result Study
	filter := bson.D{{"studyID", study.StudyID}}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.StudyUsers{Users: result.Users}, nil
}

func (s *server)  AssignUserToStudy(ctx context.Context, userStudy *pb.UserAssignment) (*pb.StudyMetaData, error) {
	var result Study
	filter := bson.D{{"studyid", userStudy.StudyID}}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	if result.Users[0] == ""{
		update1 := bson.M{"$push": bson.M{"users": userStudy.UserID}}
		updateResult1, err1 := studyCollection.UpdateOne(context.TODO(), filter, update1)
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult1.MatchedCount, updateResult1.ModifiedCount)
		update2 := bson.M{"$pop": bson.M{"users": -1}}
		updateResult2, err2 := studyCollection.UpdateOne(context.TODO(), filter, update2)
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult2.MatchedCount, updateResult2.ModifiedCount)
		if err1 != nil || err2 != nil {
			log.Fatal(err)
		}
	} else {
		update := bson.M{"$push": bson.M{"users": userStudy.UserID}}

		updateResult, err := studyCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	}
	/*var result Study
	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}*/

	return &pb.StudyMetaData{}, nil
}