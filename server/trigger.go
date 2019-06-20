package main

import (
	pb "Thesis-demo/api"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var triggerCollection = client.Database("test").Collection("triggers")
var triggerResults []*Trigger

type Trigger struct {
	Condition[] string
	StudyID string
	Action[]* pb.Action
}

type Condition struct{
	Condition[] string
}

type Attributes struct {
	age int32
	sick string
	weight int32
}

func contains(condtions []bool, search bool) bool {
	for _, value := range condtions {
		if value == search {
			return true
		}
	}
	return false
}

func createTriggerDocument(doc Trigger) {
	insertResult, err := triggerCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(insertResult)
}



func getAllTriggers() []*Trigger{
	if err != nil{
		log.Fatal(err)
	}

	findOptions := options.Find()
	cur, err := triggerCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Trigger
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		triggerResults = append(triggerResults, &elem)
	}
	return triggerResults
}


func (s *server)  CreateTrigger(ctx context.Context, createdTrigger *pb.CreatedTrigger) (*pb.CreatedTrigger, error) {

	trigger := Trigger{createdTrigger.Condition, createdTrigger.StudyID, createdTrigger.Action}
	createTriggerDocument(trigger)

	log.Printf("Trigger Created: %v", trigger.Condition)
	return &pb.CreatedTrigger{Condition: createdTrigger.Condition, StudyID: createdTrigger.StudyID,  Action: createdTrigger.Action}, nil
}

func (s *server)  CheckTrigger(attributes *pb.Attributes, streamAction pb.Trigger_CheckTriggerServer)  error {
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
