package main

import (
	pb "StudyManagement/api"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var triggerCollection = client.Database("test").Collection("triggers")
var triggerResults []*Trigger

type Trigger struct {
	ID primitive.ObjectID  `bson:"_id,omitempty"`
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



