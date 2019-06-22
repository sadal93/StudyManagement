package main

import (
	pb "Thesis-demo/api"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

var studyCollection = client.Database("test").Collection("study")
var studyResults []*Study

type Study struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
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

func getAllStudies() []*Study{
	if err != nil{
		log.Fatal(err)
	}

	findOptions := options.Find()
	cur, err := studyCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Study
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		studyResults = append(studyResults, &elem)
	}
	return studyResults
}

func (s *server)  CreateStudy(ctx context.Context, study *pb.StudyMetaData) (*pb.StudyMetaData, error) {

	studyDoc:= Study{primitive.NewObjectID(),study.Name, study.Description, time.Now().UnixNano() / 1000000, study.Status,study.Users}
	createStudyDocument(studyDoc)

	log.Printf("Study Created: %v", study.Name)
	return &pb.StudyMetaData{Name: study.Name, Description: study.Description, StartDate: study.StartDate, Status: study.Status, Users: study.Users}, nil
}

func (s *server)  GetUsers(ctx context.Context, study *pb.StudyID) (*pb.StudyUsers, error) {

	var result Study
	objectID, err := primitive.ObjectIDFromHex(study.StudyID)

	filter := bson.M{"_id": objectID}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Users for this study: %v", result.Users)

	return &pb.StudyUsers{Users: result.Users}, nil
}

func (s *server)  AssignUserToStudy(ctx context.Context, userStudy *pb.UserAssignment) (*pb.StudyMetaData, error) {
	var result Study
	objectID, err := primitive.ObjectIDFromHex(userStudy.StudyID)
	filter := bson.M{"_id": objectID}

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

	return &pb.StudyMetaData{}, nil
}

func (s *server) GetAll(ctx context.Context, empty *pb.Empty) (*pb.StudyArray, error) {

	var studies []*pb.StudyMetaData
	documents := getAllStudies()
	for _, document := range documents{
		if document.Status == "Active"{
			var study *pb.StudyMetaData = new(pb.StudyMetaData)
			study.Id = document.ID.Hex()
			study.Name = document.Name
			study.Description = document.Description
			study.Status = document.Status
			study.StartDate = document.StartDate
			study.Users = document.Users

			studies = append(studies, study)
		}
	}
	return &pb.StudyArray{Studies: studies}, nil
}

func (s *server) UpdateStudy(ctx context.Context, study *pb.StudyMetaData) (*pb.StudyMetaData, error){

	var result Study
	objectID, err := primitive.ObjectIDFromHex(study.Id)
	fmt.Println(objectID)
	filter := bson.M{"_id": objectID}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	update := bson.M{
		"$set": bson.M{
			"name": study.Name, "description": study.Description, "status": study.Status}}

	updateResult, err := studyCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return &pb.StudyMetaData{Id: study.Id, Name: study.Name, Description: study.Description, StartDate: study.StartDate, Status: study.Status, Users: study.Users}, nil
}

func (s *server) DeleteStudy(ctx context.Context, study *pb.StudyID) (*pb.Empty, error) {

	objectID, err := primitive.ObjectIDFromHex(study.StudyID)
	filter := bson.M{"_id": objectID}
	deleteResult, err := studyCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	return &pb.Empty{}, nil
}

func (s *server) GetStudy(ctx context.Context, study *pb.StudyID) (*pb.StudyMetaData, error) {
	var result Study
	objectID, err := primitive.ObjectIDFromHex(study.StudyID)
	filter := bson.M{"_id": objectID}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.StudyMetaData{Id: result.ID.Hex(), Name: result.Name, Description: result.Description, Status: result.Status, StartDate:result.StartDate, Users: result.Users }, nil

}

func (s *server) AssignWeeklySurvey(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {

	var users[] string
	documents := getAllStudies()

	for _, document := range documents{
		users = document.Users
		var survey Survey
		filter := bson.M{"study": document, "type": "weekly"}
		err = surveyCollection.FindOne(context.TODO(), filter).Decode(&survey)
		for _, user := range users{
			var result User
			filter := bson.M{"userid": user}
			err = userCollection.FindOne(context.TODO(), filter).Decode(&result)
			if (time.Now().UnixNano() / 1000000) - result.TimeLastAssigned >= result.TimeToSend{
				assignSurveyDataDoc:= AssignSurvey{primitive.NewObjectID(), survey.ID.Hex(), user, document.ID.Hex()}
				createAssignSurveyDocument(assignSurveyDataDoc)
			}
		}

	}
	time.Sleep(time.Minute)

	return &pb.Empty{},nil
}