package main

import (
	pb "Thesis-demo/api"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
	"log"
	"reflect"
	"strconv"
	"strings"
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

func (s *server) DeleteStudy(ctx context.Context, study *pb.StudyMetaData) (*pb.Empty, error) {

	objectID, err := primitive.ObjectIDFromHex(study.Id)
	filter := bson.M{"_id": objectID}
	deleteResult, err := studyCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	return &pb.Empty{}, nil
}

func (s *server) GetStudy(ctx context.Context, study *pb.StudyMetaData) (*pb.StudyMetaData, error) {
	var result Study
	objectID, err := primitive.ObjectIDFromHex(study.Id)
	filter := bson.M{"_id": objectID}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.StudyMetaData{Id: result.ID.Hex(), Name: result.Name, Description: result.Description, Status: result.Status, StartDate:result.StartDate, Users: result.Users }, nil

}

func assignWeeklySurvey() {

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
}


func (s *server)  CreateTrigger(ctx context.Context, createdTrigger *pb.CreatedTrigger) (*pb.CreatedTrigger, error) {

	trigger := Trigger{createdTrigger.Condition, createdTrigger.StudyID, createdTrigger.Action}
	createTriggerDocument(trigger)

	log.Printf("Trigger Created: %v", trigger.Condition)
	return &pb.CreatedTrigger{Condition: createdTrigger.Condition, StudyID: createdTrigger.StudyID,  Action: createdTrigger.Action}, nil
}

func (s *server)  CheckTrigger(attributes *pb.Attributes, streamAction pb.Study_CheckTriggerServer)  error {
	documents := getAllTriggers()
	//fmt.Println("USER: " , attributes)
	var actions []*pb.Action

	attributesObj := Attributes{ attributes.Age, attributes.Sick, attributes.Weight}

	for _,document := range documents {
		checks := document.Condition
		fmt.Println("Document: ", document)
		var conditions []bool

		for  _,check := range checks{
			condition := false
			var operator string
			var attribute string
			var value int64
			var valueString string
			if strings.Contains(check, "<"){
				condition := strings.Split(check, "<")
				attribute = condition[0]
				if i, err := strconv.Atoi(condition[1]); err == nil {
					fmt.Println(" i: ", i)
					value = int64(i)
					fmt.Println(" value: ", value)
				} else {
					valueString = condition[1]
					fmt.Println(" value: ", valueString)
				}
				fmt.Println(" attribute: ", attribute)

				operator = "<"
				fmt.Println(" operator: ", operator)

			} else if strings.Contains(check, ">"){
				condition := strings.Split(check, ">")
				attribute = condition[0]
				if i, err := strconv.Atoi(condition[1]); err == nil {
					fmt.Println(" i = ", i)
					value = int64(i)
					fmt.Println(" value: ", value)
				} else {
					valueString = condition[1]
					fmt.Println(" value: ", valueString)
				}

				fmt.Println(" attribute: ", attribute)

				operator = ">"
				fmt.Println(" operator: ", operator)

			} else {
				condition := strings.Split(check, "=")
				attribute = condition[0]
				if i, err := strconv.Atoi(condition[1]); err == nil {
					fmt.Println(" i: ", i)
					value = int64(i)
					fmt.Println(" value: ", value)
				} else {
					valueString = condition[1]
					fmt.Println(" value: ", valueString)
				}

				fmt.Println(" attribute: ", attribute)

				operator = "="
				fmt.Println(" operator: ", operator)
			}

			switch operator {
			case ">":
				rv := reflect.ValueOf(attributesObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if rv.FieldByName(attribute).Int() > value{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)

			case "<":
				rv := reflect.ValueOf(attributesObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if rv.FieldByName(attribute).Int() < value{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)

			case "=":
				rv := reflect.ValueOf(attributesObj)
				//rv = rv.Elem()
				fmt.Println("  rv: ", rv)
				fmt.Println("  field by name: ", rv.FieldByName(attribute))
				if  rv.FieldByName(attribute).String() == valueString{
					condition = true
				}
				fmt.Println(" CONDITION ", condition)
				conditions = append(conditions, condition)
			}

		}
		fmt.Println("CONDITIONS", conditions)
		fmt.Println("")

		if !contains(conditions, false){
			actions = document.Action
		}
	}

	for _,action := range actions{
		if action.Type == "survey"{

		}
		err := streamAction.Send(&pb.Action{Type:action.Type, Value:action.Value})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server)  UserSignUp(ctx context.Context, userStudy *pb.SignUpData) (*pb.UserMetaData, error) {

	userDoc:= User{primitive.NewObjectID(), 0, 3600000}
	createUserDocument(userDoc)

	var result Study
	objectID, err := primitive.ObjectIDFromHex(userStudy.StudyID)
	filter := bson.M{"_id": objectID}

	err = studyCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Users[0] == ""{
		update1 := bson.M{"$push": bson.M{"users": userStudy.User.Id}}
		updateResult1, err1 := studyCollection.UpdateOne(context.TODO(), filter, update1)
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult1.MatchedCount, updateResult1.ModifiedCount)

		update2 := bson.M{"$pop": bson.M{"users": -1}}
		updateResult2, err2 := studyCollection.UpdateOne(context.TODO(), filter, update2)
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult2.MatchedCount, updateResult2.ModifiedCount)

		if err1 != nil || err2 != nil {
			log.Fatal(err)
		}

	} else {
		update := bson.M{"$push": bson.M{"users": userStudy.User.Id}}

		updateResult, err := studyCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	}

	log.Printf("User Created: %v", userStudy.User.Id)
	return &pb.UserMetaData{}, nil
}